package middlewares

import (
	"cafeteller-api/firebase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthMiddleware checks if the request is authenticated and optionally checks for admin privileges
func AuthMiddleware(requireAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		userToken := c.GetHeader("Authorization")
		client := firebase.GetAuthClient(c)

		// Remove "Bearer " from the token
		userToken = strings.Replace(userToken, "Bearer ", "", 1)

		// Verify the token
		token, err := client.VerifyIDToken(c, userToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Get user ID from token
		userID := token.UID

		// Check for custom claims (e.g., isAdmin)
		if requireAdmin {
			claims := token.Claims
			if isAdmin, ok := claims["isAdmin"].(bool); !ok || !isAdmin {
				c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Admin access required"})
				c.Abort()
				return
			}
		}

		// Get user data
		_, err = client.GetUser(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
			c.Abort()
			return
		}

		// If everything is good, let the request proceed
		c.Next()
	}
}
