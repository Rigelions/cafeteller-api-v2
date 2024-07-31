package handler

import (
	"context"
	"net/http"

	"cafeteller-api/firebase"

	cloud_firestore "cloud.google.com/go/firestore"

	"github.com/gin-gonic/gin"
)

func HelloWorld(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello HHH",
	})
}

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
