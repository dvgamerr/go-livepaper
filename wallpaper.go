package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png" // imported for image.Decode
	"os"

	"github.com/nfnt/resize"
)

// loadAndResizeImage loads an image from a path and resizes it to the given dimensions #
func loadAndResizeImage(path string, width, height uint) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// but we're ensuring PNG is registered by importing image/png
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return resize.Resize(width, height, img, resize.Lanczos3), nil
	// return img, nil
}

// createBlackCanvas creates a new RGBA canvas with black background
func createBlackCanvas(width, height int) *image.RGBA {
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: image.Black}, image.Point{}, draw.Src)
	return canvas
}

// saveImageAsJPEG saves an image to the specified path as JPEG
func saveImageAsJPEG(img image.Image, path string, quality int) error {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	return jpeg.Encode(outFile, img, &jpeg.Options{Quality: quality})
}

// // mergeImages creates a black canvas and places two images on it at specified positions and sizes
// func mergeImages(fgPath1, fgPath2, outputPath string, posX1, posY1, posX2, posY2 int, fgWidth1, fgHeight1, fgWidth2, fgHeight2 uint, canvasWidth, canvasHeight int) error {
// 	canvas := createBlackCanvas(canvasWidth, canvasHeight)

// 	// Load and draw first image
// 	img1, err := loadAndResizeImage(fgPath1, fgWidth1, fgHeight1)
// 	if err != nil {
// 		return err
// 	}
// 	draw.Draw(canvas, img1.Bounds().Add(image.Pt(posX1, posY1)), img1, image.Point{}, draw.Over)

// 	// Load and draw second image
// 	img2, err := loadAndResizeImage(fgPath2, fgWidth2, fgHeight2)
// 	if err != nil {
// 		return err
// 	}
// 	draw.Draw(canvas, img2.Bounds().Add(image.Pt(posX2, posY2)), img2, image.Point{}, draw.Over)

// 	return saveImageAsJPEG(canvas, outputPath, 90)
// }
