package router

import (
	"cafeteller-api/handler"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", handler.HealthCheck)

	r.GET("/banner", handler.GetRecommendReviews)

	r.GET("/get-similar-cafe", handler.GetSimilarCafe)

	setupReviewRoute(r)
	setupAuthRoute(r)
	setupMediaRoute(r)
	setupCafeRoute(r)
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"https://tunnel.cafeteller.club", "https://cafeteller.club", "http://localhost:3000", "https://dev.cafeteller.club", "https://pre-production.cafeteller.club"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	r.Use(cors.New(config))

	limiter := tollbooth.NewLimiter(10, nil)

	// Apply the rate limiter middleware to the router
	r.Use(tollbooth_gin.LimitHandler(limiter))

	RegisterRoutes(r)
	return r
}
