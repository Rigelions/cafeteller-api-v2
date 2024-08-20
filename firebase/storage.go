package firebase

import (
	storage2 "cloud.google.com/go/storage"
	"firebase.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func GetStorageClient(c *gin.Context) *storage.Client {
	client := StorageClient
	if client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firestore client not initialized"})
		return nil
	}
	return client
}

func GetBucketName() string {
	name := os.Getenv("BUCKET_NAME")
	if name == "" {
		return ""
	}
	return name
}

func GetBucket(c *gin.Context) *storage2.BucketHandle {
	name := GetBucketName()
	if name == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bucket name not initialized"})
		return nil
	}
	client := GetStorageClient(c)
	bucket, err := client.Bucket(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bucket"})
	}

	return bucket
}
