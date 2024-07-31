package firebase

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

// FirestoreClient is a global Firestore client
var FirestoreClient *firestore.Client

// InitializeFirebase initializes the Firebase app and Firestore client
func InitializeFirebase() {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	credFilePath := os.Getenv("SERVICE_ACCOUNT_KEY_PATH")

	conf := &firebase.Config{ProjectID: projectID}
	opt := option.WithCredentialsFile(credFilePath)
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln("Error initializing Firestore client:", err)
	}
	FirestoreClient = client
}

// GetFirestoreClient retrieves the Firestore client and handles error checking
func GetFirestoreClient(c *gin.Context) *firestore.Client {
	client := FirestoreClient
	if client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firestore client not initialized"})
		return nil
	}
	return client
}

// CloseFirestoreClient closes the Firestore client
func CloseFirestoreClient() {
	if FirestoreClient != nil {
		FirestoreClient.Close()
	}
}
