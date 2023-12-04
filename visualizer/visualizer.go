package visualizer

import (
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"os"
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
func (fi *FractalImage) WriteImage() {
	// Save the image to a file
	outputFile, err := os.Create(fi.Path)
	if err != nil {
		slog.Error("Error creating file:", err)
		return
	}
	defer outputFile.Close()
	png.Encode(outputFile, fi.Img)
}

// will "thicken" the image to become a 3d object
func (fi *FractalImage) ExtrudeImage(start, end float32) [][]float32{
	var mesh [][]float32

	for y := fi.Img.Bounds().Min.Y; y < fi.Img.Bounds().Max.Y-1; y++ {
		for x := fi.Img.Bounds().Min.X; x < fi.Img.Bounds().Max.X-1; x++ {
			// Get the color and alpha values of the current pixel
			_, _, _, a := fi.Img.At(x, y).RGBA()

			// Check if the pixel is not fully transparent
			if a > 0 {
				// Define the vertices of a quad
				v0 := []float32{float32(x), float32(y), start}
				v1 := []float32{float32(x + 1), float32(y), start}
				v2 := []float32{float32(x + 1), float32(y + 1), start}
				v3 := []float32{float32(x), float32(y + 1), start}

				// Extrude the quad to the specified thickness
				v4 := []float32{float32(x), float32(y), end}
				v5 := []float32{float32(x + 1), float32(y), end}
				v6 := []float32{float32(x + 1), float32(y + 1), end}
				v7 := []float32{float32(x), float32(y + 1), end}

				// Add the two triangles (facets) of the extruded quad to the mesh
				mesh = append(mesh, v0, v1, v2, v3, v4, v5)
				mesh = append(mesh, v2, v3, v6, v7, v4, v5)
			}
		}
	}

	return mesh
}