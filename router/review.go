package router

import (
	"cafeteller-api/handler"
	"github.com/gin-gonic/gin"
)

func setupReviewRoute(r *gin.Engine) {
	r.GET("/reviews", handler.GetReviews)
	r.GET("/reviews/:id", handler.GetReviewByID)
}
