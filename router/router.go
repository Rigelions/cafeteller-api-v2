package router

import (
	"cafeteller-api/handler"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
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

	r.GET("/hello", handler.HelloWorld)

	r.GET("/reviews/:id", handler.GetReviewByID)
	return r
}
