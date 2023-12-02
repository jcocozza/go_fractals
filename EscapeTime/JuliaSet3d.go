package EscapeTime

import (
	"fmt"
	"image"
	"os"
)

const (
	cuboidWidth  = 1.0
	cuboidHeight = 1.0
	cuboidDepth  = 10.1
)

func DrawJuliaSet3D(img image.Image, stlFile *os.File, shift float64) {
	// Loop through pixels and generate cuboids for non-transparent pixels
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			_, _, _, alpha := pixel.RGBA()

			if alpha > 0 {
				writeCuboid(stlFile, float64(x), float64(y), shift)
			}
		}
	}
}

func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func writeCuboid(file *os.File, x, y, z float64) {
	// Write ASCII STL representation of a cuboid
	// Bottom face
	writeFacet(file, x, y, z, x+cuboidWidth, y, z, x+cuboidWidth, y+cuboidHeight, z)
	writeFacet(file, x, y, z, x+cuboidWidth, y+cuboidHeight, z, x, y+cuboidHeight, z)

	// Top face
	writeFacet(file, x, y, z+cuboidDepth, x+cuboidWidth, y, z+cuboidDepth, x+cuboidWidth, y+cuboidHeight, z+cuboidDepth)
	writeFacet(file, x, y, z+cuboidDepth, x+cuboidWidth, y+cuboidHeight, z+cuboidDepth, x, y+cuboidHeight, z+cuboidDepth)

	// Side faces
	writeFacet(file, x, y, z, x, y+cuboidHeight, z, x, y, z+cuboidDepth)
	writeFacet(file, x+cuboidWidth, y, z, x+cuboidWidth, y, z+cuboidDepth, x+cuboidWidth, y+cuboidHeight, z)
	writeFacet(file, x+cuboidWidth, y+cuboidHeight, z, x+cuboidWidth, y+cuboidHeight, z+cuboidDepth, x, y+cuboidHeight, z)
	writeFacet(file, x, y+cuboidHeight, z, x, y+cuboidHeight, z+cuboidDepth, x, y, z)
}

func writeFacet(file *os.File, x1, y1, z1, x2, y2, z2, x3, y3, z3 float64) {
	// Write ASCII STL representation of a facet
	file.WriteString(fmt.Sprintf("  facet normal 0 0 0\n"))
	file.WriteString(fmt.Sprintf("    outer loop\n"))
	file.WriteString(fmt.Sprintf("      vertex %f %f %f\n", x1, y1, z1))
	file.WriteString(fmt.Sprintf("      vertex %f %f %f\n", x2, y2, z2))
	file.WriteString(fmt.Sprintf("      vertex %f %f %f\n", x3, y3, z3))
	file.WriteString(fmt.Sprintf("    endloop\n"))
	file.WriteString(fmt.Sprintf("  endfacet\n"))
}