package image

import (
	"time"
)

// ImageMetadata holds metadata about an image
type ImageMetadata struct {
	ID          string
	Width       int
	Height      int
	CapturedAt  time.Time
	Rating      int // 1-5
	Filters     []string
	Format      string // "JPEG", "PNG", etc.
	Description string
}

// Image represents a photo with its data and metadata
type Image interface {
	ID() string
	Data() []byte
	Metadata() ImageMetadata
	SetData([]byte)
	SetMetadata(ImageMetadata)
}

// BasicImage is a concrete implementation of Image
type BasicImage struct {
	id       string
	data     []byte
	metadata ImageMetadata
}

// NewBasicImage creates a new BasicImage
func NewBasicImage(id string, data []byte, metadata ImageMetadata) *BasicImage {
	metadata.ID = id
	return &BasicImage{
		id:       id,
		data:     data,
		metadata: metadata,
	}
}

// ID returns the image ID
func (b *BasicImage) ID() string {
	return b.id
}

// Data returns the image data
func (b *BasicImage) Data() []byte {
	return b.data
}

// Metadata returns the image metadata
func (b *BasicImage) Metadata() ImageMetadata {
	return b.metadata
}

// SetData updates the image data
func (b *BasicImage) SetData(data []byte) {
	b.data = data
}

// SetMetadata updates the image metadata
func (b *BasicImage) SetMetadata(metadata ImageMetadata) {
	b.metadata = metadata
}
