package main

import (
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

func main() {

	monitors := getMonitors()

	log.Printf("Total monitors: %d\n", len(monitors))

	for _, monitor := range monitors {
		primaryStatus := " "
		if monitor.primary {
			primaryStatus = "*"
		}
		log.Printf("Monitor %d%s: %+v\n", monitor.index+1, primaryStatus, monitor.resolution)
	}

	// Defaulting to Span style as per original behavior
	if err := setWallpaperStyle(STYLE_SPAN); err != nil {
		log.Printf("Error setting wallpaper style: %v\n", err)
	}

	// IMPORTANT: Use double backslashes in Go string literals for Windows paths
	imagePath := "D:\\home\\Downloads\\twitter\\Drowsy_sheep\\Drowsy_sheep-1910560294415315392-01.jpg"
	if err := setWallpaper(imagePath); err != nil {
		log.Printf("Error setting wallpaper: %v\n", err)
		os.Exit(1)
	}
}
