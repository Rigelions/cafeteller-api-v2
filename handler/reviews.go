package handler

import (
	"cafeteller-api/firebase"
	"context"
	"net/http"
	"time"

	cloud_firestore "cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func GetReviewByID(c *gin.Context) {
	ctx := context.Background()
	id := c.Param("id")

	// Use Firestore client
	client := firebase.GetFirestoreClient(c)

	dsnap, err := client.Collection("reviews").Doc(id).Get(ctx)

	if err != nil {
		//  show bad request not found
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Review not found",
		})
		return
	}

	data := dsnap.Data()

	cafe_snap, err := data["cafe"].(*cloud_firestore.DocumentRef).Get(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cafe not found",
		})
		return
	}

	data["cafe"] = cafe_snap.Data()

	c.JSON(http.StatusOK, data)
}

func GetReviews(c *gin.Context) {
	ctx := context.Background()

	// Use Firestore client
	client := firebase.GetFirestoreClient(c)

	// if c has last_id, get reviews after that id
	before_update := c.Query("update_date")

	var query cloud_firestore.Query

	if before_update != "" {
		// Convert before_update to the appropriate type, e.g., a Firestore timestamp
		parsedBeforeUpdate, err := time.Parse(time.RFC3339, before_update)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid before_update format"})
			return
		}

		query = client.Collection("reviews").OrderBy("updateDate", cloud_firestore.Desc).StartAfter(parsedBeforeUpdate).Limit(9)
	} else {
		query = client.Collection("reviews").OrderBy("updateDate", cloud_firestore.Desc).Limit(9)
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	var reviews []map[string]interface{}

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		reviewData := doc.Data()
		reviewData["id"] = doc.Ref.ID

		reviews = append(reviews, reviewData)
	}

	// get cafes that has reviews in the list of doc ref
	reviewsRefs := make([]*cloud_firestore.DocumentRef, len(reviews))
	reviewMap := make(map[string]map[string]interface{})

	for i, review := range reviews {
		reviewsRefs[i] = client.Collection("reviews").Doc(review["id"].(string))
		reviewMap[review["id"].(string)] = review
	}

	cafeQuery := client.Collection("cafes").Where("reviews", "in", reviewsRefs)

	cafeIter := cafeQuery.Documents(ctx)
	defer cafeIter.Stop()

	for {
		doc, err := cafeIter.Next()
		if err != nil {
			break
		}

		cafe_data := doc.Data()
		review_ref_cafe := cafe_data["reviews"].(*cloud_firestore.DocumentRef)

		reviewID := review_ref_cafe.ID
		if review, exists := reviewMap[reviewID]; exists {
			review["cafe"] = cafe_data
		}

		// copy createDate and updateDate to from reviews to cafe
		cafe_data["createDate"] = reviewMap[reviewID]["createDate"]
		cafe_data["updateDate"] = reviewMap[reviewID]["updateDate"]
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
	})
}

func GetRecommendReviews(c *gin.Context) {
	ctx := context.Background()

	// Use Firestore client
	client := firebase.GetFirestoreClient(c)

	query := client.Collection("cafes").OrderBy("updateDate", cloud_firestore.Desc).Where("tags", "array-contains", "Recommended").Limit(3)

	queryLatest := client.Collection("cafes").OrderBy("updateDate", cloud_firestore.Desc).Limit(2)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var recommendCafes []map[string]interface{}
	var latestCafes []map[string]interface{}

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		cafeData := doc.Data()
		cafeData["id"] = doc.Ref.ID
		cafeData["review_id"] = cafeData["reviews"].(*cloud_firestore.DocumentRef).ID

		recommendCafes = append(recommendCafes, cafeData)
	}

	iterLatest := queryLatest.Documents(ctx)
	defer iterLatest.Stop()

	for {
		doc, err := iterLatest.Next()
		if err != nil {
			break
		}

		cafeData := doc.Data()
		cafeData["id"] = doc.Ref.ID
		cafeData["review_id"] = cafeData["reviews"].(*cloud_firestore.DocumentRef).ID

		latestCafes = append(latestCafes, cafeData)
	}

	c.JSON(http.StatusOK, gin.H{
		"recommend_cafes": recommendCafes,
		"latest_cafes":    latestCafes,
	})
}
