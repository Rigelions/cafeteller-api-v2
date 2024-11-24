package handler

import (
	"bytes"
	"cafeteller-api/firebase"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/gin-gonic/gin"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"image/png"
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

// ResizeImage resizes an image to the specified width while maintaining the aspect ratio.
// Returns the resized image.
func ResizeImage(img image.Image, width int) image.Image {
	bounds := img.Bounds()
	height := bounds.Dy() * width / bounds.Dx() // Maintain aspect ratio
	resizedImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, bounds, draw.Over, nil)
	return resizedImg
}

// EncodeImageToWebP encodes an image to WebP format.
// Returns the WebP image buffer or an error.
func EncodeImageToWebP(img image.Image) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	err := webp.Encode(&buffer, img, &webp.Options{Lossless: false, Quality: 100})
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}
	return &buffer, nil
}

// HandleUploadAndConvertToWebP handles the upload and conversion of an image to WebP.
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

	// Define sizes for WebP conversion
	sizes := []int{1980, 1024, 720}
	urls := make(map[string]string)
	unix := time.Now().Unix()

	// Process each size
	for _, size := range sizes {
		// Resize the image
		resizedImg := ResizeImage(img, size)

		// Skip WebP encoding; work with the original or resized image buffer
		buffer := new(bytes.Buffer)

		// Encode the resized image to its original format (e.g., JPEG or PNG)
		switch header.Header.Get("Content-Type") {
		case "image/jpeg":
			err = jpeg.Encode(buffer, resizedImg, nil) // Use JPEG encoder
		case "image/png":
			err = png.Encode(buffer, resizedImg) // Use PNG encoder
		default:
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error_message": "Unsupported image format"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
			return
		}

		// Create a unique filename for each size, maintaining the original extension
		filename := fmt.Sprintf("%d/%s@%d%s", unix, strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename)), size, filepath.Ext(header.Filename))

		// Upload to storage with the original MIME type
		url, err := uploadToStorage(c, buffer, header.Header.Get("Content-Type"), filename, "", "images")
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
