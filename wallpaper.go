package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png" // imported for image.Decode
	"os"

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

func saveImageAsJPEG(img image.Image, path string, quality int) error {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	return jpeg.Encode(outFile, img, &jpeg.Options{Quality: quality})
}
