package main

import (
	"cafeteller-api/firebase"
	"cafeteller-api/router"

	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Initialize Firebase
	firebase.InitializeFirebase()
	defer firebase.CloseFirestoreClient()

	r := router.SetupRouter()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
