package firebase

import (
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetFirestoreClient retrieves the Firestore client and handles error checking
func GetFirestoreClient(c *gin.Context) *firestore.Client {
	client := FirestoreClient
	if client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firestore client not initialized"})
		return nil
	}
	return client
}
