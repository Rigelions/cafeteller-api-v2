package router

import (
	"cafeteller-api/handler"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
)

func RegisterRoutes(api *gin.RouterGroup) {
	api.GET("/health", handler.HealthCheck)

	api.GET("/banner", handler.GetRecommendReviews)

	api.GET("/get-similar-cafe", handler.GetSimilarCafe)

	setupReviewRoute(api)
	setupAuthRoute(api)
	setupMediaRoute(api)
	setupCafeRoute(api)
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"https://tunnel.cafeteller.club", "http://localhost:3000", "https://dev.cafeteller.club"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	r.Use(cors.New(config))

	limiter := tollbooth.NewLimiter(10, nil)

	// Apply the rate limiter middleware to the router
	r.Use(tollbooth_gin.LimitHandler(limiter))

	// Add /api prefix
	api := r.Group("/api")
	RegisterRoutes(api)
	return r
}
