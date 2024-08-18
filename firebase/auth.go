package firebase

import (
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAuthClient(c *gin.Context) *auth.Client {
	client := AuthClient
	if client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth client not initialized"})
		return nil
	}

	return client
}
