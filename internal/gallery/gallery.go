// Package gallery implements the Strategy pattern for flexible image sorting.
// Different sorting algorithms (by date, rating, ID) can be swapped at runtime
// without modifying the Gallery class.
package gallery

import (
	"sort"

	"photoapp/internal/image"
)

// Gallery holds a collection of images (Context)
type Gallery struct {
	images []image.Image
	sorter Sorter
}

// NewGallery creates a new gallery
func NewGallery() *Gallery {
	return &Gallery{
		images: make([]image.Image, 0),
	}
}

// AddImage adds an image to the gallery
func (g *Gallery) AddImage(img image.Image) {
	g.images = append(g.images, img)
}

// Images returns all images
func (g *Gallery) Images() []image.Image {
	return g.images
}

// SetSorter sets the sorting strategy
func (g *Gallery) SetSorter(sorter Sorter) {
	g.sorter = sorter
}

// Sort sorts the gallery using the current strategy
func (g *Gallery) Sort() {
	if g.sorter != nil {
		g.images = g.sorter.Sort(g.images)
	}
}

// Sorter defines the strategy interface (Strategy pattern)
type Sorter interface {
	Sort(images []image.Image) []image.Image
	Name() string
}

// SortByDate sorts by date (Concrete Strategy)
type SortByDate struct {
	ascending bool
}

func NewSortByDate(ascending bool) *SortByDate {
	return &SortByDate{ascending: ascending}
}

func (s *SortByDate) Sort(images []image.Image) []image.Image {
	sorted := make([]image.Image, len(images))
	copy(sorted, images)
	sort.Slice(sorted, func(i, j int) bool {
		if s.ascending {
			return sorted[i].Metadata().CapturedAt.Before(sorted[j].Metadata().CapturedAt)
		}
		return sorted[i].Metadata().CapturedAt.After(sorted[j].Metadata().CapturedAt)
	})
	return sorted
}

func (s *SortByDate) Name() string {
	if s.ascending {
		return "Date(Asc)"
	}
	return "Date(Desc)"
}

// SortByRating sorts by rating (Concrete Strategy)
type SortByRating struct {
	ascending bool
}

func NewSortByRating(ascending bool) *SortByRating {
	return &SortByRating{ascending: ascending}
}

func (s *SortByRating) Sort(images []image.Image) []image.Image {
	sorted := make([]image.Image, len(images))
	copy(sorted, images)
	sort.Slice(sorted, func(i, j int) bool {
		if s.ascending {
			return sorted[i].Metadata().Rating < sorted[j].Metadata().Rating
		}
		return sorted[i].Metadata().Rating > sorted[j].Metadata().Rating
	})
	return sorted
}

func (s *SortByRating) Name() string {
	if s.ascending {
		return "Rating(Asc)"
	}
	return "Rating(Desc)"
}

// SortByID sorts by ID (Concrete Strategy)
type SortByID struct {
	ascending bool
}

func NewSortByID(ascending bool) *SortByID {
	return &SortByID{ascending: ascending}
}

func (s *SortByID) Sort(images []image.Image) []image.Image {
	sorted := make([]image.Image, len(images))
	copy(sorted, images)
	sort.Slice(sorted, func(i, j int) bool {
		if s.ascending {
			return sorted[i].ID() < sorted[j].ID()
		}
		return sorted[i].ID() > sorted[j].ID()
	})
	return sorted
}

func (s *SortByID) Name() string {
	if s.ascending {
		return "ID(Asc)"
	}
	return "ID(Desc)"
}
