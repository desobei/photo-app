package events

import "fmt"

const defaultThumbnailSize = 128

// LoggerObserver logs events to stdout.
type LoggerObserver struct {
	name string
}

// NewLoggerObserver creates a new logger observer.
func NewLoggerObserver(name string) *LoggerObserver {
	if name == "" {
		name = "Logger"
	}
	return &LoggerObserver{name: name}
}

// OnEvent handles incoming events by logging them.
func (l *LoggerObserver) OnEvent(event *Event) {
	if event == nil {
		return
	}
	fmt.Printf("[%s] %s\n", l.name, event.Type)
}

// Name returns the observer name.
func (l *LoggerObserver) Name() string {
	return l.name
}

// ThumbnailGeneratorObserver generates image thumbnails.
type ThumbnailGeneratorObserver struct {
	name       string
	thumbnails map[string][]byte
}

// NewThumbnailGeneratorObserver creates a new thumbnail generator.
func NewThumbnailGeneratorObserver(name string) *ThumbnailGeneratorObserver {
	if name == "" {
		name = "ThumbnailGenerator"
	}
	return &ThumbnailGeneratorObserver{
		name:       name,
		thumbnails: make(map[string][]byte),
	}
}

// OnEvent handles events by generating thumbnails.
func (t *ThumbnailGeneratorObserver) OnEvent(event *Event) {
	if event == nil || event.Image == nil {
		return
	}
	data := event.Image.Data()
	size := min(len(data), defaultThumbnailSize)
	thumb := make([]byte, size)
	copy(thumb, data[:size])
	t.thumbnails[event.Image.ID()] = thumb
}

// Name returns the observer name.
func (t *ThumbnailGeneratorObserver) Name() string {
	return t.name
}

// GetThumbnail retrieves a thumbnail by image ID.
func (t *ThumbnailGeneratorObserver) GetThumbnail(imageID string) ([]byte, bool) {
	thumb, ok := t.thumbnails[imageID]
	return thumb, ok
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// StatisticsObserver tracks event statistics.
type StatisticsObserver struct {
	name  string
	count map[EventType]int
}

// NewStatisticsObserver creates a new statistics observer.
func NewStatisticsObserver(name string) *StatisticsObserver {
	if name == "" {
		name = "Statistics"
	}
	return &StatisticsObserver{
		name:  name,
		count: make(map[EventType]int),
	}
}

// OnEvent handles events by incrementing counters.
func (s *StatisticsObserver) OnEvent(event *Event) {
	if event == nil {
		return
	}
	s.count[event.Type]++
}

// Name returns the observer name.
func (s *StatisticsObserver) Name() string {
	return s.name
}

// GetStats returns formatted event statistics.
func (s *StatisticsObserver) GetStats() string {
	result := "Statistics:\n"
	for eventType, count := range s.count {
		result += fmt.Sprintf("  %s: %d\n", eventType, count)
	}
	return result
}
