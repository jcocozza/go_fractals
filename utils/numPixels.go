package utils

import "image"

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