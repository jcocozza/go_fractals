package IteratedFunctionSystems

import (
	"image"
	"image/color"

	IMGS "github.com/jcocozza/go_fractals/images"
)

type FractalImage struct {
	Width int
	Height int
	Path string
	PointsList [][]float64
	Img *image.RGBA
}

//NewFractalImage will create a new image
func NewFractalImage(width int, height int, path string, pointsList [][]float64) *FractalImage {
	return &FractalImage{
		Width: width,
		Height: height,
		Path: path,
		PointsList: pointsList,
		Img: _DrawFractal(width,height,pointsList),
	}
}

func _DrawFractal(width int, height int, pointsList [][]float64) *image.RGBA {
	img := image.NewRGBA(image.Rect(0,0,width,height))
	minX, minY, maxX, maxY := pointsList[0][0], pointsList[0][1], pointsList[0][0], pointsList[0][1]

	// Find the range of coordinates among all generated points
	for _, point := range pointsList {
		x, y := point[0], point[1]
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	// Iterate over collected points and draw them on the canvas
	for i := range pointsList {
		point := pointsList[i]
		// Draw a point on the canvas
		scaledX := int((point[0] - minX) / (maxX - minX) * float64(width))
		scaledY := int((point[1] - minY) / (maxY - minY) * float64(height))

		// Draw a point on the canvas
		img.Set(scaledX, height - scaledY, color.RGBA{0, 255, 0, 255})
		}

	return img
}

// write the image of a fractal to a png
func (fi *FractalImage) WriteImage(path string) {
	// Save the image to a file
	IMGS.SavePNG(fi.Img, path)
}