package handler

import (
	"cafeteller-api/firebase"
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Cafe struct {
	Name                     string      `json:"name"`
	Location                 Location    `json:"location"`
	Tags                     []string    `json:"tags"`
	SublocalityLevel1        string      `json:"sublocality_level_1"`
	AdministrativeAreaLevel1 string      `json:"administrative_area_level_1"`
	CreateDate               interface{} `json:"createDate"`
	UpdateDate               interface{} `json:"updateDate"`
}

type Location struct {
	Latitude  float64 `json:"lat" firestore:"lat"`
	Longitude float64 `json:"lon" firestore:"lon"`
}

func AutocompleteCafeName(c *gin.Context) {
	query := c.Query("name")

	client := firebase.GetFirestoreClient(c)

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	// Call Firestore to search for matching cafes
	cafes, err := searchCafes(client, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//
	c.JSON(http.StatusOK, cafes)
}

func searchCafes(client *firestore.Client, query string) ([]Cafe, error) {
	var cafes []Cafe
	ctx := context.Background()

	// Firestore query to search for cafes with name matching the query (case-insensitive)
	docs, err := client.Collection("cafes").Where("name_search", ">=", query).Where("name_search", "<=", query+"\uf8ff").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		var cafe Cafe
		if err := doc.DataTo(&cafe); err != nil {
			return nil, err
		}
		cafes = append(cafes, cafe)
	}

	return cafes, nil
}

func MigrateCafeNamesToLowercase(c *gin.Context) {
	ctx := context.Background()
	client := firebase.GetFirestoreClient(c)

	// Fetch all cafes
	docs, err := client.Collection("cafes").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cafe documents"})
		return
	}

	// Get Reviews from cafe
	reviews, err := client.Collection("reviews").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews documents"})
		return
	}

	// map review.cafe.id with cafe
	cafeMap := make(map[string]map[string]interface{})
	for _, review := range reviews {
		reviewData := review.Data()
		reviewData["id"] = review.Ref.ID
		cafeRef := reviewData["cafe"].(*firestore.DocumentRef)
		cafeMap[cafeRef.ID] = reviewData
	}

	// Loop through each document and update the name to lowercase
	for _, doc := range docs {
		var cafe Cafe
		if err := doc.DataTo(&cafe); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse cafe document: " + doc.Ref.ID})
			return
		}

		// copy createDate and updateDate to from reviews to cafe
		if review, exists := cafeMap[doc.Ref.ID]; exists {
			cafe.CreateDate = review["createDate"]
			cafe.UpdateDate = review["updateDate"]
		}

		// Check if the name is already lowercase
		lowercaseName := strings.ToLower(cafe.Name)
		if cafe.Name != lowercaseName {
			// Update the name to lowercase
			_, err = client.Collection("cafes").Doc(doc.Ref.ID).Update(ctx, []firestore.Update{
				{Path: "name_search", Value: lowercaseName},
				{Path: "createDate", Value: cafe.CreateDate},
				{Path: "updateDate", Value: cafe.UpdateDate},
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cafe: " + cafe.Name})
				return
			}
		}
	}

	// If everything is successful
	c.JSON(http.StatusOK, gin.H{"message": "All cafe names migrated to lowercase successfully!"})
}

func SearchCafesWithFiltersHandler(c *gin.Context) {
	client := firebase.GetFirestoreClient(c)

	// Extract parameters from the query string
	query := c.DefaultQuery("name", "")
	amphoe := c.QueryArray("amphoe")     // Receives amphoe as a list
	province := c.QueryArray("province") // Receives amphoe as a list
	tags := c.QueryArray("tags")         // Receives tags as a list

	// Call the search function
	cafes, err := searchCafesWithFilters(client, query, amphoe, tags, province)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the result as JSON
	c.JSON(http.StatusOK, cafes)
}

func searchCafesWithFilters(client *firestore.Client, name string, amphoes []string, tags []string, province []string) ([]Cafe, error) {
	var cafes []Cafe
	ctx := context.Background()

	// Start Firestore query from the "cafes" collection
	collection := client.Collection("cafes")
	firestoreQuery := collection.Query

	// Apply name query (if provided)
	if name != "" {
		firestoreQuery = firestoreQuery.Where("name_search", "==", name)
	}

	// Apply amphoe filter (if provided)
	if len(amphoes) > 0 {
		firestoreQuery = firestoreQuery.Where("sublocality_level_1", "in", amphoes)
	}

	fmt.Println(len(province))
	// Apply province filter (if provided)
	if len(province) > 0 {
		firestoreQuery = firestoreQuery.Where("administrative_area_level_1", "in", province)
	}

	// Apply tags filter (if provided)
	if len(tags) > 0 {
		firestoreQuery = firestoreQuery.Where("tags", "array-contains-any", tags)
	}

	// Limit to 10 results if no filters are applied
	if name == "" && len(amphoes) == 0 && len(tags) == 0 && len(province) == 0 {
		firestoreQuery = firestoreQuery.Limit(10)
	}

	// Execute the query
	docs, err := firestoreQuery.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	// Map Firestore documents to Cafe structs
	for _, doc := range docs {
		var cafe Cafe
		if err := doc.DataTo(&cafe); err != nil {
			return nil, err
		}
		cafes = append(cafes, cafe)
	}

	return cafes, nil
}
