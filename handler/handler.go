package handler

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"

	"cafeteller-api/firebase"

	cloudFirestore "cloud.google.com/go/firestore"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

func HelloWorld(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello test live",
	})
}

// GetSimilarCafe handler function
func GetSimilarCafe(c *gin.Context) {
	ctx := context.Background()
	client := firebase.GetFirestoreClient(c)

	// Get tags from query parameters
	tags := c.QueryArray("tags[]")

	if len(tags) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tags provided"})
		return
	}

	// Create a query to find cafes with matching tags
	query := client.Collection("cafes").Where("tags", "array-contains-any", tags)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var cafes []map[string]interface{}
	cafeMap := make(map[string]map[string]interface{})

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching cafes"})
			return
		}

		cafeData := doc.Data()
		cafeData["id"] = doc.Ref.ID

		// Calculate similarity score based on tags intersection
		similarityScore := calculateSimilarityScore(tags, cafeData["tags"].([]interface{}))
		cafeData["similarityScore"] = similarityScore

		cafes = append(cafes, cafeData)
		cafeMap[doc.Ref.ID] = cafeData
	}

	// Sort cafes by similarity score in descending order
	sort.SliceStable(cafes, func(i, j int) bool {
		return cafes[i]["similarityScore"].(int) > cafes[j]["similarityScore"].(int)
	})

	// Extract cafe IDs for further querying
	cafeIDs := extractCafeIDs(cafes)

	// Get reviews for the similar cafes
	reviews, err := getReviewsForCafes(ctx, client, cafeIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching reviews"})
		return
	}

	// Map the reviews with their corresponding cafes
	for _, review := range reviews {
		if cafe, exists := cafeMap[review["cafe"].(*cloudFirestore.DocumentRef).ID]; exists {
			review["cafe"] = cafe
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
	})
}

// Function to calculate similarity score
func calculateSimilarityScore(requestTags []string, cafeTags []interface{}) int {
	score := 0
	for _, tag := range requestTags {
		for _, cafeTag := range cafeTags {
			if strings.EqualFold(tag, cafeTag.(string)) {
				score++
			}
		}
	}
	return score
}

// Function to extract cafe IDs
func extractCafeIDs(cafes []map[string]interface{}) []string {
	ids := make([]string, len(cafes))
	for i, cafe := range cafes {
		ids[i] = cafe["id"].(string)
	}
	return ids
}

// Function to get reviews for specific cafes
func getReviewsForCafes(ctx context.Context, client *cloudFirestore.Client, cafeIDs []string) ([]map[string]interface{}, error) {
	reviews := []map[string]interface{}{}
	var cafeRefs []*cloudFirestore.DocumentRef

	for _, id := range cafeIDs {
		cafeRefs = append(cafeRefs, client.Collection("cafes").Doc(id))
	}

	// Create a query to find reviews for the cafes
	query := client.Collection("reviews").Where("cafe", "in", cafeRefs)

	iter := query.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		reviewData := doc.Data()
		reviewData["id"] = doc.Ref.ID
		reviews = append(reviews, reviewData)
	}

	return reviews, nil
}
