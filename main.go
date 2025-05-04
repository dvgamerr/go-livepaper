package main

import (
	"fmt"
	"os"
)

func main() {

	monitors := getMonitors()

	fmt.Printf("Total monitors: %d\n\n", len(monitors))

	for _, monitor := range monitors {
		primaryStatus := ""
		if monitor.Primary {
			primaryStatus = " (Primary)"
		}
		fmt.Printf("Monitor %d%s:\n", monitor.Index+1, primaryStatus)
		fmt.Printf("  Resolution: %s\n", monitor.Resolution)
		fmt.Printf("  Position: %s\n", monitor.Position)
	}

	// IMPORTANT: Use double backslashes in Go string literals for Windows paths
	imagePath := "D:\\home\\Downloads\\twitter\\Drowsy_sheep\\Drowsy_sheep-1910560294415315392-01.jpg"

	fmt.Printf("Setting wallpaper to: %s\n", imagePath)
	if err := setWallpaper(imagePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting wallpaper: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Wallpaper set successfully!")
}
