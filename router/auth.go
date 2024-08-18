package router

import (
	"cafeteller-api/handler"
	"github.com/gin-gonic/gin"
)

func setupAuthRoute(r *gin.Engine) {
	r.POST("/auth", handler.Auth)
}
