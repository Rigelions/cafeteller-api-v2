package router

import (
	"cafeteller-api/handler"
	"cafeteller-api/middlewares"
	"github.com/gin-gonic/gin"
)

func setupAuthRoute(r *gin.Engine) {
	r.GET("/auth", handler.Auth)
	r.POST("/set-admin", middlewares.AuthMiddleware(true), handler.SetAdmin)
}
