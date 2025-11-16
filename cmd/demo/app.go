package demo

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"photoapp/internal/camera"
	"photoapp/internal/codec"
	"photoapp/internal/events"
	"photoapp/internal/gallery"
	"photoapp/internal/image"
	"photoapp/internal/storage"
)

type App struct {
	eventBus  *events.EventBus
	facade    *camera.Facade
	gallery   *gallery.Gallery
	storage   storage.Storage
	loggerObs *events.LoggerObserver
	thumbObs  *events.ThumbnailGeneratorObserver
	statsObs  *events.StatisticsObserver
	scanner   *bufio.Scanner
}

func NewApp() *App {
	eventBus := events.NewEventBus()

	// Register observers
	loggerObs := events.NewLoggerObserver("SystemLogger")
	thumbObs := events.NewThumbnailGeneratorObserver("ThumbnailGen")
	statsObs := events.NewStatisticsObserver("StatsTracker")

	eventBus.Register(loggerObs)
	eventBus.Register(thumbObs)
	eventBus.Register(statsObs)

	// Create storage adapter (Adapter pattern)
	store := storage.NewMapAdapter()

	facade := camera.NewFacade(eventBus, store)
	gal := gallery.NewGallery()

	return &App{
		eventBus:  eventBus,
		facade:    facade,
		gallery:   gal,
		storage:   store,
		loggerObs: loggerObs,
		thumbObs:  thumbObs,
		statsObs:  statsObs,
		scanner:   bufio.NewScanner(os.Stdin),
	}
}

func Run() {
	app := NewApp()
	app.showWelcome()
	app.mainMenu()
}

func (a *App) showWelcome() {
	fmt.Println("   Photo Gallery App - Design Patterns Demonstration      ")
	fmt.Println()
}

func (a *App) mainMenu() {
	for {
		fmt.Println("\n┌─────────────────── MAIN MENU ───────────────────┐")
		fmt.Println("│ 1. Capture Photo                                │")
		fmt.Println("│ 2. View Gallery                                 │")
		fmt.Println("│ 3. Sort Gallery                                 │")
		fmt.Println("│ 4. Demo Factory Patterns                        │")
		fmt.Println("│ 5. Demo Decorator Pattern                       │")
		fmt.Println("│ 6. View Statistics                              │")
		fmt.Println("│ 7. View Thumbnails                              │")
		fmt.Println("│ 0. Exit                                         │")
		fmt.Println("└─────────────────────────────────────────────────┘")
		fmt.Print("Select option: ")

		choice := a.readInput()
		fmt.Println()

		switch choice {
		case "1":
			a.capturePhoto()
		case "2":
			a.viewGallery()
		case "3":
			a.sortGallery()
		case "4":
			a.demoFactories()
		case "5":
			a.demoDecorator()
		case "6":
			a.viewStatistics()
		case "7":
			a.viewThumbnails()
		case "0":
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid option. Please try again.")
		}
	}
}

func (a *App) capturePhoto() {
	fmt.Println("📷 Capture Photo")
	fmt.Println("─────────────────")

	// Select photo type
	fmt.Println("\nSelect photo type:")
	fmt.Println("1. Landscape")
	fmt.Println("2. Portrait")
	fmt.Print("Choice: ")
	photoChoice := a.readInput()

	photoType := "landscape"
	if photoChoice == "2" {
		photoType = "portrait"
	}

	// Select filters
	fmt.Println("\nSelect filters (comma-separated, or press Enter for none):")
	fmt.Println("Available: grayscale, sepia, blur")
	fmt.Print("Filters: ")
	filtersInput := a.readInput()

	var filters []string
	if filtersInput != "" {
		filters = strings.Split(strings.ReplaceAll(filtersInput, " ", ""), ",")
	}

	// Select format
	fmt.Println("\nSelect format:")
	fmt.Println("1. JPEG")
	fmt.Println("2. PNG")
	fmt.Print("Choice: ")
	formatChoice := a.readInput()

	format := "jpeg"
	if formatChoice == "2" {
		format = "png"
	}

	fmt.Printf("\n📸 Creating %s photo, filters=%v, format=%s...\n", photoType, filters, format)

	encoded, err := a.facade.CaptureAndProcess(photoType, filters, format)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}

	// Add to gallery
	img, _ := a.facade.QuickCapture(photoType)

	// Apply filters if any
	if len(filters) > 0 {
		processedImg := image.Image(img)
		for _, filter := range filters {
			processedImg = image.NewFilterDecorator(processedImg, filter)
		}
		img = processedImg
	}

	a.gallery.AddImage(img)

	fmt.Printf("✅ Photo created! Size: %d bytes\n", len(encoded))
	fmt.Printf("   Image ID: %s\n", img.ID())
	fmt.Printf("   Total photos in gallery: %d\n", len(a.gallery.Images()))
}

func (a *App) viewGallery() {
	fmt.Println("🖼️  Gallery")
	fmt.Println("─────────────────")

	images := a.gallery.Images()
	if len(images) == 0 {
		fmt.Println("📭 Gallery is empty. Capture some photos first!")
		return
	}

	fmt.Printf("Total images: %d\n\n", len(images))
	for i, img := range images {
		meta := img.Metadata()
		fmt.Printf("[%d] ID: %s\n", i+1, img.ID())
		fmt.Printf("    Filters: %v\n", meta.Filters)
		fmt.Printf("    Rating: %d\n", meta.Rating)
		fmt.Printf("    Format: %s\n", meta.Format)
		fmt.Printf("    Captured: %s\n", meta.CapturedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("    Size: %d bytes\n\n", len(img.Data()))
	}
}

func (a *App) sortGallery() {
	fmt.Println("📊 Sort Gallery")
	fmt.Println("─────────────────")

	if len(a.gallery.Images()) == 0 {
		fmt.Println("📭 Gallery is empty. Nothing to sort!")
		return
	}

	fmt.Println("\nSelect sorting strategy:")
	fmt.Println("1. Sort by Date (Ascending)")
	fmt.Println("2. Sort by Date (Descending)")
	fmt.Println("3. Sort by Rating (Ascending)")
	fmt.Println("4. Sort by Rating (Descending)")
	fmt.Print("Choice: ")

	choice := a.readInput()

	var sorter gallery.Sorter
	var sortName string

	switch choice {
	case "1":
		sorter = gallery.NewSortByDate(true)
		sortName = "Date (Ascending)"
	case "2":
		sorter = gallery.NewSortByDate(false)
		sortName = "Date (Descending)"
	case "3":
		sorter = gallery.NewSortByRating(true)
		sortName = "Rating (Ascending)"
	case "4":
		sorter = gallery.NewSortByRating(false)
		sortName = "Rating (Descending)"
	default:
		fmt.Println("❌ Invalid choice")
		return
	}

	a.gallery.SetSorter(sorter)
	a.gallery.Sort()

	a.eventBus.Notify(events.NewEvent(
		events.EventGallerySorted,
		nil,
		fmt.Sprintf("Gallery sorted by: %s", sortName),
	))

	fmt.Printf("✅ Gallery sorted by: %s\n", sortName)
}

func (a *App) demoFactories() {
	fmt.Println("🏭 Factory + Adapter Patterns Demo")
	fmt.Println("─────────────────")

	// Factory Method
	fmt.Println("\n📌 Factory Method Pattern - Creating Photos")
	photoFactory := camera.NewFactory()

	landscape := photoFactory.CreatePhoto("landscape")
	fmt.Printf("  ✓ Created: %s\n", landscape.ID())

	portrait := photoFactory.CreatePhoto("portrait")
	fmt.Printf("  ✓ Created: %s\n", portrait.ID())

	// Codec (simple encoder/decoder)
	fmt.Println("\n📌 Codec - Image Encoders/Decoders")

	jpegEncoder := codec.NewJPEGEncoder()
	fmt.Printf("\n  JPEG Encoder: %s\n", jpegEncoder.Format())
	jpegDecoder := codec.NewJPEGDecoder()
	fmt.Printf("  JPEG Decoder: %s\n", jpegDecoder.Format())

	pngEncoder := codec.NewPNGEncoder()
	fmt.Printf("\n  PNG Encoder: %s\n", pngEncoder.Format())
	pngDecoder := codec.NewPNGDecoder()
	fmt.Printf("  PNG Decoder: %s\n", pngDecoder.Format())

	// Adapter Pattern
	fmt.Println("\n📌 Adapter Pattern - Map Storage Adapter")
	fmt.Printf("  Current storage: Map Adapter\n")
	fmt.Printf("  Storage adapter adapts map to Storage interface\n")
}

func (a *App) demoDecorator() {
	fmt.Println("🎨 Decorator Pattern Demo")
	fmt.Println("─────────────────")

	fmt.Println("\n📸 Creating base photo...")
	baseImg, err := a.facade.QuickCapture("portrait")
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}

	fmt.Printf("Base photo: %d bytes, filters=%v\n", len(baseImg.Data()), baseImg.Metadata().Filters)

	// Stack decorators
	fmt.Println("\n🔧 Applying decorator chain: Grayscale → Sepia → Blur")
	processedImg := image.Image(baseImg)
	processedImg = image.NewFilterDecorator(processedImg, "Grayscale")
	processedImg = image.NewFilterDecorator(processedImg, "Sepia")
	processedImg = image.NewFilterDecorator(processedImg, "Blur")

	fmt.Printf("\n✅ Processing complete!\n")
	fmt.Printf("   Original filters: %v\n", baseImg.Metadata().Filters)
	fmt.Printf("   After decorators: %v\n", processedImg.Metadata().Filters)
	fmt.Printf("   Data size: %d bytes (unchanged)\n", len(processedImg.Data()))
}

func (a *App) viewStatistics() {
	fmt.Println("📈 Statistics")
	fmt.Println("─────────────────")
	fmt.Println(a.statsObs.GetStats())
	fmt.Printf("\nGallery size: %d images\n", len(a.gallery.Images()))
	fmt.Printf("Observers registered: %d\n", a.eventBus.ObserverCount())
}

func (a *App) viewThumbnails() {
	fmt.Println("🖼️  Thumbnails")
	fmt.Println("─────────────────")

	images := a.gallery.Images()
	if len(images) == 0 {
		fmt.Println("📭 No thumbnails available. Capture some photos first!")
		return
	}

	count := 0
	for _, img := range images {
		if thumb, ok := a.thumbObs.GetThumbnail(img.ID()); ok {
			fmt.Printf("  • %s: %d bytes\n", img.ID(), len(thumb))
			count++
		}
	}

	fmt.Printf("\nTotal thumbnails: %d\n", count)
}

func (a *App) readInput() string {
	a.scanner.Scan()
	return strings.TrimSpace(a.scanner.Text())
}

func (a *App) readInt() int {
	input := a.readInput()
	val, _ := strconv.Atoi(input)
	return val
}
