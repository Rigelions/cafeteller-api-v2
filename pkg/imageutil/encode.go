package imageutil

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/chai2010/webp"
)

func DecodeImage(file io.Reader, format string) (image.Image, error) {
	switch format {
	case "jpeg":
		return jpeg.Decode(file)
	case "png":
		return png.Decode(file)
	case "webp":
		return webp.Decode(file)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}
}

func EncodeImage(buffer *bytes.Buffer, img image.Image, format string) error {
	switch format {
	case "jpeg":
		return jpeg.Encode(buffer, img, nil)
	case "png":
		return png.Encode(buffer, img)
	case "webp":
		return webp.Encode(buffer, img, &webp.Options{Lossless: false, Quality: 80})
	default:
		return fmt.Errorf("unsupported image format: %s", format)
	}
}
