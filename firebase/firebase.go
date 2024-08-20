package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/storage"
	"google.golang.org/api/option"
	"log"
	"os"
)

// FirestoreClient is a global Firestore client
var FirestoreClient *firestore.Client
var AuthClient *auth.Client
var StorageClient *storage.Client

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

	aClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalln("Error initializing Auth client:", err)
	}
	AuthClient = aClient

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln("Error initializing Firestore client:", err)
	}
	FirestoreClient = client

	storageClient, err := app.Storage(ctx)
	if err != nil {
		log.Fatalln("Error initializing Storage client:", err)
	}
	StorageClient = storageClient
}

// CloseFirestoreClient closes the Firestore client
func CloseFirestoreClient() {
	if FirestoreClient != nil {
		err := FirestoreClient.Close()
		if err != nil {
			return
		}
	}
}
