package imageutil

import (
	"github.com/nfnt/resize"
	"image"
)

func ResizeImage(img image.Image, maxSize int) image.Image {
	// Get original dimensions
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	// No resizing if both dimensions are already smaller than maxSize
	if originalWidth <= maxSize && originalHeight <= maxSize {
		return img
	}

	// Calculate the new dimensions
	var newWidth, newHeight uint
	if originalWidth > originalHeight {
		newWidth = uint(maxSize)
		newHeight = uint((maxSize * originalHeight) / originalWidth)
	} else {
		newHeight = uint(maxSize)
		newWidth = uint((maxSize * originalWidth) / originalHeight)
	}

	// Resize the image
	return resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
}
