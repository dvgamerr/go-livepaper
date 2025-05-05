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
	"github.com/rwcarlsen/goexif/exif"
)

func resizeImageToFill(img image.Image, orientation int, width, height uint) (image.Image, image.Point) {
	img = applyOrientation(img, orientation)

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

func applyOrientation(img image.Image, orientation int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	switch orientation {
	case 1:
		// 0 degrees: the correct orientation, no adjustment required
		return img
	case 2:
		// 0 degrees, mirrored: image has been flipped back-to-front
		dst := image.NewRGBA(bounds)
		for y := range height {
			for x := range width {
				dst.Set(width-x-1, y, img.At(x, y))
			}
		}
		return dst
	case 3:
		// 180 degrees: image is upside down
		dst := image.NewRGBA(bounds)
		for y := range height {
			for x := range width {
				dst.Set(width-x-1, height-y-1, img.At(x, y))
			}
		}
		return dst
	case 4:
		// 180 degrees, mirrored: image has been flipped back-to-front and is upside down
		dst := image.NewRGBA(bounds)
		for y := range height {
			for x := range width {
				dst.Set(x, height-y-1, img.At(x, y))
			}
		}
		return dst
	case 5:
		// 90 degrees: image has been flipped back-to-front and is on its side
		dst := image.NewRGBA(image.Rect(0, 0, height, width))
		for y := range height {
			for x := range width {
				dst.Set(y, width-x-1, img.At(x, y))
			}
		}
		return dst
	case 6:
		// 90 degrees, mirrored: image is on its side
		dst := image.NewRGBA(image.Rect(0, 0, height, width))
		for y := range height {
			for x := range width {
				dst.Set(height-y-1, x, img.At(x, y))
			}
		}
		return dst
	case 7:
		// 270 degrees: image has been flipped back-to-front and is on its far side
		dst := image.NewRGBA(image.Rect(0, 0, height, width))
		for y := range height {
			for x := range width {
				dst.Set(y, x, img.At(x, y))
			}
		}
		return dst
	case 8:
		// 270 degrees, mirrored: image is on its far side
		dst := image.NewRGBA(image.Rect(0, 0, height, width))
		for y := range height {
			for x := range width {
				dst.Set(height-y-1, width-x-1, img.At(x, y))
			}
		}
		return dst
	default:
		// Unknown orientation, return original image
		return img
	}
}

func loadAndResizeImage(path string, width, height uint) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the image file
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	// Rewind file to beginning to read EXIF data
	if _, err = file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek file: %w", err)
	}

	// Try to extract EXIF data
	orientation := 1
	exifData, err := exif.Decode(file)
	if err != nil {
		// Continue even if EXIF extraction fails
		fmt.Println("Warning: Could not extract EXIF data:", err)
	} else {
		// Process EXIF data if needed
		// For example, to check orientation:
		if data, err := exifData.Get(exif.Orientation); err == nil {
			if val, err := data.Int(0); err == nil {
				orientation = val
				fmt.Println("Image orientation:", val)
			}
		}
	}

	resized, offset := resizeImageToFill(img, orientation, width, height)
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
