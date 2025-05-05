package main

import (
	"fmt"
	"image"
	"image/draw"
	"log"

	"github.com/alexflint/go-arg"
)

// VERSION is set during build using ldflags
// go build -ldflags "-X main.VERSION=1.0.0"
var VERSION = "dev"

func (Args) Version() string {
	return fmt.Sprintf("go-livepaper %s", VERSION)
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

// Args defines the command line arguments
type Args struct {
	Monitor   []string `arg:"-m,--monitor" help:"Target monitors to set wallpaper (e.g. -m 1 -m 2)"`
	Wallpaper []string `arg:"positional" help:"Wallpaper is a list of file paths to wallpaper images"`
}

var args Args

func main() {
	// Define command line arguments
	arg.MustParse(&args)
	if len(args.Wallpaper) == 0 {
		log.Fatalf("No wallpaper specified. Please provide at least one wallpaper image path.")
	}
	if len(args.Wallpaper) != len(args.Monitor) && len(args.Monitor) > 0 {
		log.Fatalf("Invalid arguments: monitors (%d) must match of wallpapers (%d)", len(args.Monitor), len(args.Wallpaper))
	}

	canvasWidth, canvasHeight, monitors := getMonitors()
	log.Printf("  Monitor: %d\n", len(monitors))
	log.Printf("Wallpaper: %dx%dpx", canvasWidth, canvasHeight)

	// Defaulting to Span style as per original behavior
	if err := setWallpaperStyle(STYLE_SPAN); err != nil {
		log.Printf("Error setting wallpaper style: %v\n", err)
	}

	canvas := createBlackCanvas(canvasWidth, canvasHeight)

	for i, monitor := range monitors {

		primaryStatus := " "
		if monitor.primary {
			primaryStatus = "*"

		}
		log.Printf("Monitor %d%s: %+v\n", monitor.index+1, primaryStatus, monitor.resolution)

		if i > len(args.Wallpaper)-1 {
			continue
		}
		img1, err := loadAndResizeImage(args.Wallpaper[i], uint(monitor.resolution.width), uint(monitor.resolution.height))
		if err != nil {
			log.Printf("Error loadAndResizeImage: %v\n", err)
		}
		draw.Draw(canvas, img1.Bounds().Add(image.Pt(int(monitor.resolution.x), int(monitor.resolution.y))), img1, image.Point{}, draw.Over)

	}

	// // Load and draw first image
	// img1, err := loadAndResizeImage(fgPath1, fgWidth1, fgHeight1)
	// if err != nil {
	// 	return err
	// }
	// draw.Draw(canvas, img1.Bounds().Add(image.Pt(posX1, posY1)), img1, image.Point{}, draw.Over)

	// // Load and draw second image
	// img2, err := loadAndResizeImage(fgPath2, fgWidth2, fgHeight2)
	// if err != nil {
	// 	return err
	// }
	// draw.Draw(canvas, img2.Bounds().Add(image.Pt(posX2, posY2)), img2, image.Point{}, draw.Over)

	saveImageAsJPEG(canvas, "./output/test.jpg", 100)

	// // IMPORTANT: Use double backslashes in Go string literals for Windows paths
	// imagePath := "D:\\home\\Downloads\\twitter\\Drowsy_sheep\\Drowsy_sheep-1910560294415315392-01.jpg"
	// if err := setWallpaper(imagePath); err != nil {
	// 	log.Printf("Error setting wallpaper: %v\n", err)
	// 	os.Exit(1)
	// }
}
