package router

import (
	"cafeteller-api/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/hello", handler.HelloWorld)
	return r
}
