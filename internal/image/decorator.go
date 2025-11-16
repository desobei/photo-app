// Package image implements the Decorator pattern.
package image

// FilterDecorator wraps an Image and adds filter metadata.
type FilterDecorator struct {
	wrapped Image
	filter  string
}

// NewFilterDecorator creates a new filter decorator.
func NewFilterDecorator(img Image, filter string) *FilterDecorator {
	if img == nil {
		panic("image cannot be nil")
	}
	if filter == "" {
		panic("filter cannot be empty")
	}
	return &FilterDecorator{wrapped: img, filter: filter}
}

func (d *FilterDecorator) ID() string {
	return d.wrapped.ID()
}

func (d *FilterDecorator) Data() []byte {
	return d.wrapped.Data()
}

func (d *FilterDecorator) Metadata() ImageMetadata {
	meta := d.wrapped.Metadata()
	meta.Filters = append(meta.Filters, d.filter)
	return meta
}

func (d *FilterDecorator) SetData(data []byte) {
	d.wrapped.SetData(data)
}

func (d *FilterDecorator) SetMetadata(meta ImageMetadata) {
	d.wrapped.SetMetadata(meta)
}
