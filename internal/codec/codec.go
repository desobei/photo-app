// Package codec provides image encoding and decoding functionality.
// Simple encoder/decoder implementations for JPEG and PNG formats.
package codec

import (
	"fmt"

	"photoapp/internal/image"
)

// Encoder encodes an image to bytes
type Encoder interface {
	Encode(img image.Image) ([]byte, error)
	Format() string
}

// Decoder decodes bytes to an image
type Decoder interface {
	Decode(data []byte) (image.Image, error)
	Format() string
}

// JPEGEncoder encodes images to JPEG format
type JPEGEncoder struct{}

// NewJPEGEncoder creates a new JPEG encoder
func NewJPEGEncoder() *JPEGEncoder {
	return &JPEGEncoder{}
}

// Encode simulates JPEG encoding
func (e *JPEGEncoder) Encode(img image.Image) ([]byte, error) {
	data := img.Data()
	// Simulate JPEG compression (prefix with magic bytes)
	encoded := append([]byte{0xFF, 0xD8, 0xFF, 0xE0}, data...)
	return encoded, nil
}

// Format returns the format name
func (e *JPEGEncoder) Format() string {
	return "JPEG"
}

// JPEGDecoder decodes JPEG images
type JPEGDecoder struct{}

// NewJPEGDecoder creates a new JPEG decoder
func NewJPEGDecoder() *JPEGDecoder {
	return &JPEGDecoder{}
}

// Decode simulates JPEG decoding
func (d *JPEGDecoder) Decode(data []byte) (image.Image, error) {
	if len(data) < 4 || data[0] != 0xFF || data[1] != 0xD8 {
		return nil, fmt.Errorf("invalid JPEG data")
	}
	// Remove JPEG header
	rawData := data[4:]

	metadata := image.ImageMetadata{
		Format: "JPEG",
		Width:  1920,
		Height: 1080,
	}

	return image.NewBasicImage("decoded-jpeg", rawData, metadata), nil
}

// Format returns the format name
func (d *JPEGDecoder) Format() string {
	return "JPEG"
}

// PNGEncoder encodes images to PNG format
type PNGEncoder struct{}

// NewPNGEncoder creates a new PNG encoder
func NewPNGEncoder() *PNGEncoder {
	return &PNGEncoder{}
}

// Encode simulates PNG encoding
func (e *PNGEncoder) Encode(img image.Image) ([]byte, error) {
	data := img.Data()
	// Simulate PNG encoding (prefix with PNG signature)
	encoded := append([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, data...)
	return encoded, nil
}

// Format returns the format name
func (e *PNGEncoder) Format() string {
	return "PNG"
}

// PNGDecoder decodes PNG images
type PNGDecoder struct{}

// NewPNGDecoder creates a new PNG decoder
func NewPNGDecoder() *PNGDecoder {
	return &PNGDecoder{}
}

// Decode simulates PNG decoding
func (d *PNGDecoder) Decode(data []byte) (image.Image, error) {
	if len(data) < 8 || data[0] != 0x89 || data[1] != 0x50 {
		return nil, fmt.Errorf("invalid PNG data")
	}
	// Remove PNG header
	rawData := data[8:]

	metadata := image.ImageMetadata{
		Format: "PNG",
		Width:  1920,
		Height: 1080,
	}

	return image.NewBasicImage("decoded-png", rawData, metadata), nil
}

// Format returns the format name
func (d *PNGDecoder) Format() string {
	return "PNG"
}
