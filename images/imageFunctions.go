package images

import (
	"image"
	"image/png"
	"log/slog"
	"os"
)

// save to a png
func SavePNG(img image.Image, path string) {
	outfile, err := os.Create(path)
	if err != nil {
		slog.Error("Error creating image file:", err)
		return
	}
	defer outfile.Close()

	err = png.Encode(outfile, img)
	if err != nil {
		slog.Error("Error encoding image to PNG:", err)
		return
	}
}

// function to get alpha value of a pixel, handling out-of-bounds gracefully
func GetAlpha(img image.Image, x, y int) uint32 {
	bounds := img.Bounds()
	if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
		pixel := img.At(x, y)
		_, _, _, alpha := pixel.RGBA()
		return alpha
	}
	return 0
}

// if any adjacent points are transparent, then the non-transparent pixel is on the boundary
func IsTransition(img image.Image, x,y int) bool {
	alphaLeft := GetAlpha(img, x-1, y)
	alphaRight := GetAlpha(img, x+1, y)
	alphaUp := GetAlpha(img, x, y-1)
	alphaDown := GetAlpha(img, x, y+1)

	return alphaLeft == 0 || alphaRight == 0 || alphaUp == 0 || alphaDown == 0
}

// check if the pixel in the current image will be "covered" by pixels in other images.
// the specific use case is in an ordered list of images.
// x,y will be apoint in the current image, and images will be a list of future images
func IsCovered(x,y int, images []image.Image) bool {
	for _, img := range images {
		pixel := img.At(x, y)
			_, _, _, alpha := pixel.RGBA()
			if alpha > 0 {
				return true
			}
	}
	return false
}

func GetNumNonTransparentPixels(img image.Image) int {
	bounds := img.Bounds()

	// Initialize a counter for non-transparent pixels
	nonTransparentPixels := 0

	// Iterate through each pixel in the image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the color of the current pixel
			pixelColor := img.At(x, y)

			// Check if the pixel is not fully transparent
			_, _, _, alpha := pixelColor.RGBA()
			if alpha != 0 {
				nonTransparentPixels++
			}
		}
	}
	return nonTransparentPixels
}