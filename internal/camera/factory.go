// Package camera implements the Factory Method pattern.
package camera

import (
	"fmt"
	"math/rand"
	"time"

	"photoapp/internal/image"
)

const (
	defaultWidth    = 1920
	defaultHeight   = 1080
	defaultDataSize = 1024
	minRating       = 1
	maxRating       = 5
	defaultFormat   = "JPEG"
)

const (
	PhotoTypeLandscape = "landscape"
	PhotoTypePortrait  = "portrait"
)

// Factory creates photos using the Factory Method pattern.
type Factory struct{}

// NewFactory creates a new photo factory.
func NewFactory() *Factory {
	return &Factory{}
}

// CreatePhoto creates a photo of the specified type.
func (f *Factory) CreatePhoto(photoType string) image.Image {
	data := make([]byte, defaultDataSize)
	rand.Read(data)

	metadata := image.ImageMetadata{
		Width:       defaultWidth,
		Height:      defaultHeight,
		CapturedAt:  time.Now(),
		Rating:      rand.Intn(maxRating) + minRating,
		Filters:     []string{},
		Format:      defaultFormat,
		Description: f.getDescription(photoType),
	}

	id := fmt.Sprintf("photo-%d", time.Now().UnixNano())
	return image.NewBasicImage(id, data, metadata)
}

func (f *Factory) getDescription(photoType string) string {
	switch photoType {
	case PhotoTypeLandscape:
		return "Beautiful landscape photo"
	case PhotoTypePortrait:
		return "Portrait photo"
	default:
		return "Standard photo"
	}
}
