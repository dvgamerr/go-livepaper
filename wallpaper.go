package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png" // imported for image.Decode
	"os"
	"path/filepath"
	"time"

	"github.com/nfnt/resize"
)

func resizeImageToFill(img image.Image, width, height uint) (image.Image, image.Point) {
	// Get original dimensions
	bounds := img.Bounds()
	origWidth := uint(bounds.Dx())
	origHeight := uint(bounds.Dy())

	// Calculate aspect ratios
	origRatio := float64(origWidth) / float64(origHeight)
	targetRatio := float64(width) / float64(height)

	// Determine new dimensions that maintain aspect ratio and fill the target area
	var newWidth, newHeight uint
	if origRatio > targetRatio {
		// Image is wider than target, height becomes the limiting factor
		newHeight = height
		newWidth = uint(float64(height) * origRatio)
	} else {
		// Image is taller than target, width becomes the limiting factor
		newWidth = width
		newHeight = uint(float64(width) / origRatio)
	}

	// Resize the image to the new dimensions
	resized := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	// Calculate centering offsets
	offsetX := (int(newWidth) - int(width)) / 2
	offsetY := (int(newHeight) - int(height)) / 2

	return resized, image.Point{X: offsetX, Y: offsetY}
}

func loadAndResizeImage(path string, width, height uint) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	resized, offset := resizeImageToFill(img, width, height)
	result := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	draw.Draw(result, result.Bounds(), resized, offset, draw.Src)

	return result, nil
}

func createBlackCanvas(width, height int) *image.RGBA {
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: image.Black}, image.Point{}, draw.Src)
	return canvas
}

func saveImageAs(img image.Image, quality int) (string, error) {
	// Check if the path is absolute
	filename := fmt.Sprintf("wallpaper_%d.jpg", time.Now().UnixNano())

	if !filepath.IsAbs(filename) {
		// If not absolute, save to temp directory
		tempDir, err := getTempDir()
		if err != nil {
			return "", fmt.Errorf("failed to get temp directory: %w", err)
		}

		// Make sure the directory exists
		dir := filepath.Dir(filepath.Join(tempDir, filename))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}

		filename = filepath.Join(tempDir, filename)
	}

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return filename, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Encode the image to JPEG format
	err = jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return filename, fmt.Errorf("failed to encode image: %w", err)
	}

	return filename, nil
}

// getTempDir returns the application's temporary directory path, creating it if needed
func getTempDir() (string, error) {
	tempDir := filepath.Join(os.TempDir(), "go-livepaper")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	return tempDir, nil
}

// cleanTempDir removes all files in the temporary directory
func cleanTempDir() error {
	tempDir, err := getTempDir()
	if err != nil {
		return err
	}

	// Check if directory exists before attempting to remove files
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return nil // Directory doesn't exist, nothing to clean
	}

	// Remove all files in the directory
	dirEntries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %w", err)
	}

	for _, entry := range dirEntries {
		if err := os.RemoveAll(filepath.Join(tempDir, entry.Name())); err != nil {
			return fmt.Errorf("failed to remove %s: %w", entry.Name(), err)
		}
	}

	return nil
}
