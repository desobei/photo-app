// Package camera implements the Facade pattern.
package camera

import (
	"fmt"
	"strings"

	"photoapp/internal/codec"
	"photoapp/internal/events"
	"photoapp/internal/image"
	"photoapp/internal/storage"
)

const (
	FormatPNG  = "png"
	FormatJPEG = "jpeg"
)

// Facade simplifies complex photo processing workflows.
type Facade struct {
	factory  *Factory
	eventBus events.Subject
	storage  storage.Storage
}

// NewFacade creates a new Facade with required dependencies.
func NewFacade(eventBus events.Subject, store storage.Storage) *Facade {
	if eventBus == nil || store == nil {
		panic("eventBus and storage cannot be nil")
	}
	return &Facade{
		factory:  NewFactory(),
		eventBus: eventBus,
		storage:  store,
	}
}

// CaptureAndProcess creates, filters, encodes, and stores a photo.
func (f *Facade) CaptureAndProcess(photoType string, filters []string, format string) ([]byte, error) {
	photo := f.createPhoto(photoType)
	processed := f.applyFilters(photo, filters)
	encoded, err := f.encodePhoto(processed, format)
	if err != nil {
		return nil, fmt.Errorf("encode photo: %w", err)
	}

	if err := f.storage.Save(processed.ID(), encoded); err != nil {
		return nil, fmt.Errorf("save photo: %w", err)
	}

	f.eventBus.Notify(events.NewEvent(events.EventImageProcessed, processed, "Processed"))
	return encoded, nil
}

// QuickCapture creates a photo without processing.
func (f *Facade) QuickCapture(photoType string) (image.Image, error) {
	return f.createPhoto(photoType), nil
}

func (f *Facade) createPhoto(photoType string) image.Image {
	photo := f.factory.CreatePhoto(photoType)
	f.eventBus.Notify(events.NewEvent(events.EventImageCaptured, photo, "Photo created"))
	return photo
}

func (f *Facade) applyFilters(photo image.Image, filters []string) image.Image {
	processed := photo
	for _, filter := range filters {
		processed = image.NewFilterDecorator(processed, filter)
	}
	return processed
}

func (f *Facade) encodePhoto(img image.Image, format string) ([]byte, error) {
	encoder := f.selectEncoder(format)
	return encoder.Encode(img)
}

func (f *Facade) selectEncoder(format string) codec.Encoder {
	if strings.ToLower(format) == FormatPNG {
		return codec.NewPNGEncoder()
	}
	return codec.NewJPEGEncoder()
}
