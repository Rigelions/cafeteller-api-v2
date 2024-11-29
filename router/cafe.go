package router

import (
	"cafeteller-api/handler"
	"github.com/gin-gonic/gin"
)

func setupCafeRoute(r *gin.RouterGroup) {
	r.GET("/cafe", handler.GetReviews)
	r.POST("/cafe/migrate", handler.MigrateCafeNamesToLowercase)
}
