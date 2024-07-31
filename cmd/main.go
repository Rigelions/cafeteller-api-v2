package main

import (
	"cafeteller-api/router"
	"os"
)

func main() {
	r := router.SetupRouter()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
