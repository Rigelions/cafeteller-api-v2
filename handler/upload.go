package handler

import (
	"bytes"
	"cafeteller-api/firebase"
	"cafeteller-api/pkg/imageutil"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/gin-gonic/gin"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func HandleUploadURL(c *gin.Context) {
	url := c.PostForm("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "URL is required"})
		return
	}

	id := time.Now().Unix()
	// Create a unique filename
	filename := fmt.Sprintf("%d", id)

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Failed to download file"})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Failed to close response body"})
		}
	}(resp.Body)

	// Determine file type
	contentType := resp.Header.Get("Content-Type")
	filetype := "." + strings.Split(contentType, "/")[1]

	// Upload the file to Google Cloud Storage
	publicURL, err := uploadToStorage(c, resp.Body, contentType, filename, filetype, "instagram")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"success": 1,
		"id":      id,
		"file": gin.H{
			"url":      publicURL,
			"filename": filepath.Base(publicURL),
		},
	})
}

func HandleUploadAndConvertToWebP(c *gin.Context) {
	// Get the file from the form data
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "File is required"})
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Failed to close file"})
		}
	}(file)

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "Invalid image format"})
		return
	}

	// Create different sized WebP images
	sizes := []int{1980, 1024, 720}
	urls := make(map[string]string)
	unix := time.Now().Unix()

	for _, size := range sizes {
		resizedImg := imageutil.ResizeImage(img, size)
		var buffer bytes.Buffer

		// Encode the resized image to WebP
		err := webp.Encode(&buffer, resizedImg, &webp.Options{Lossless: false, Quality: 100})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Failed to encode image"})
			return
		}

		// Create a unique filename for each size
		filename := fmt.Sprintf("%d/%s@%d.webp", unix, strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename)), size)
		url, err := uploadToStorage(c, &buffer, "image/webp", filename, "", "images")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
			return
		}
		urls[fmt.Sprintf("@%d", size)] = url
	}

	// Return response with URLs of different sizes
	c.JSON(http.StatusOK, gin.H{
		"success": 1,
		"id":      unix,
		"urls":    urls,
	})
}

func uploadToStorage(c *gin.Context, file io.Reader, contentType, filename, filetype string, baseFolder string) (string, error) {
	// Set up Google Cloud Storage
	ctx := context.Background()

	// Define the path
	path := fmt.Sprintf("%s/%s/%s%s", baseFolder, time.Now().Format("20060102"), filename, filetype)
	bucket := firebase.GetBucket(c)

	object := bucket.Object(path)
	wc := object.NewWriter(ctx)
	wc.ContentType = contentType

	// Upload the file
	if _, err := io.Copy(wc, file); err != nil {
		_ = wc.Close()
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Close the writer
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close storage writer: %w", err)
	}

	// Make the file public
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("failed to make file public: %w", err)
	}

	// Get the bucket name
	bucketName := firebase.GetBucketName()

	// Construct the public URL
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, path)

	return publicURL, nil
}
