package router

import (
	"cafeteller-api/handler"
	"cafeteller-api/middlewares"
	"github.com/gin-gonic/gin"
	"os"
)

func setupAuthRoute(r *gin.RouterGroup) {
	env := os.Getenv("GO_ENV")

	if env == "development" {
		r.GET("/auth/local", handler.AuthLocal)
	}

	r.GET("/auth", handler.AuthRemote)
	r.POST("/set-admin", middlewares.AuthMiddleware(true), handler.SetAdmin)
}
