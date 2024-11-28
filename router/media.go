package router

import (
	"cafeteller-api/handler"
	"cafeteller-api/middlewares"
	"github.com/gin-gonic/gin"
)

func setupMediaRoute(r *gin.Engine) {
	r.POST("/media/ig/upload", middlewares.AuthMiddleware(true), handler.HandleUploadURL)
	r.POST("/media/image/upload", middlewares.AuthMiddleware(true), handler.HandleUploadAndConvertToWebP)
}
