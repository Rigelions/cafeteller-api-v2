package router

import (
	"cafeteller-api/handler"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
	}
	r.Use(cors.New(config))

	limiter := tollbooth.NewLimiter(10, nil)

	// Apply the rate limiter middleware to the router
	r.Use(tollbooth_gin.LimitHandler(limiter))

	r.GET("/hello", handler.HelloWorld)

	r.GET("/reviews/:id", handler.GetReviewByID)

	r.GET("/get-similar-cafe", handler.GetSimilarCafe)
	return r
}
