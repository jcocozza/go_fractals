package EscapeTime

import (
	"fmt"
	"image"
	"os"
)

const (
	cubeSize    = 1.0
	facetNormal = "  facet normal %f %f %f\n"
	outerLoop   = "    outer loop\n"
	vertex      = "      vertex %f %f %f\n"
	endLoop     = "    endloop\n"
	endFacet    = "  endfacet\n"
)

// for each non-transparent pixel in the image, draw
func DrawJuliaSet3D(img image.Image, stlFile *os.File, shift float64) {
	// Loop through pixels and generate cuboids for non-transparent pixels
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			_, _, _, alpha := pixel.RGBA()

			if alpha > 0 {
				writeCube(stlFile, float64(x), float64(y), shift)
			}
		}
	}
}

// draw a cube at an (x,y,z) coordinate
func writeCube(file *os.File, x, y, z float64) {
		// Write ASCII STL representation of a facet
		file.WriteString(fmt.Sprintf(facetNormal, 0.0, 0.0, -1.0))
		file.WriteString(fmt.Sprintf(outerLoop))
		file.WriteString(fmt.Sprintf(vertex, x, y, z))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y + cubeSize,z))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y,z))
		file.WriteString(fmt.Sprintf(endLoop))
		file.WriteString(fmt.Sprintf(endFacet))

		file.WriteString(fmt.Sprintf(facetNormal, 0.0, 0.0, -1.0))
		file.WriteString(fmt.Sprintf(outerLoop))
		file.WriteString(fmt.Sprintf(vertex, x, y, z))
		file.WriteString(fmt.Sprintf(vertex, x, y + cubeSize, z))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y + cubeSize, z))
		file.WriteString(fmt.Sprintf(endLoop))
		file.WriteString(fmt.Sprintf(endFacet))

		file.WriteString(fmt.Sprintf(facetNormal, -1.0,0.0,0.0))
		file.WriteString(fmt.Sprintf(outerLoop))
		file.WriteString(fmt.Sprintf(vertex, x, y, z))
		file.WriteString(fmt.Sprintf(vertex, x, y + cubeSize, z + cubeSize))
		file.WriteString(fmt.Sprintf(vertex, x, y + cubeSize, z))
		file.WriteString(fmt.Sprintf(endLoop))
		file.WriteString(fmt.Sprintf(endFacet))

		file.WriteString(fmt.Sprintf(facetNormal, -1.0,0.0,0.0))
		file.WriteString(fmt.Sprintf(outerLoop))
		file.WriteString(fmt.Sprintf(vertex, x, y, z))
		file.WriteString(fmt.Sprintf(vertex, x, y, z + cubeSize))
		file.WriteString(fmt.Sprintf(vertex, x, y + cubeSize, z + cubeSize))
		file.WriteString(fmt.Sprintf(endLoop))
		file.WriteString(fmt.Sprintf(endFacet))

		file.WriteString(fmt.Sprintf(facetNormal, 0.0,1.0,0.0))
		file.WriteString(fmt.Sprintf(outerLoop))
		file.WriteString(fmt.Sprintf(vertex, x, y + cubeSize, z))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y + cubeSize, z + cubeSize))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y + cubeSize, z))
		file.WriteString(fmt.Sprintf(endLoop))
		file.WriteString(fmt.Sprintf(endFacet))

		file.WriteString(fmt.Sprintf(facetNormal, 0.0,1.0,0.0))
		file.WriteString(fmt.Sprintf(outerLoop))
		file.WriteString(fmt.Sprintf(vertex, x, y + cubeSize, z))
		file.WriteString(fmt.Sprintf(vertex, x, y + cubeSize, z + cubeSize))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y + cubeSize, z + cubeSize))
		file.WriteString(fmt.Sprintf(endLoop))
		file.WriteString(fmt.Sprintf(endFacet))

		file.WriteString(fmt.Sprintf(facetNormal, 1.0,0.0,0.0))
		file.WriteString(fmt.Sprintf(outerLoop))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y, z))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y + cubeSize, z))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y + cubeSize, z + cubeSize))
		file.WriteString(fmt.Sprintf(endLoop))
		file.WriteString(fmt.Sprintf(endFacet))

		file.WriteString(fmt.Sprintf(facetNormal, 1.0,0.0,0.0))
		file.WriteString(fmt.Sprintf(outerLoop))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y, z))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y + cubeSize, z + cubeSize))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y, z + cubeSize))
		file.WriteString(fmt.Sprintf(endLoop))
		file.WriteString(fmt.Sprintf(endFacet))

		file.WriteString(fmt.Sprintf(facetNormal, 0.0,-1.0,0.0))
		file.WriteString(fmt.Sprintf(outerLoop))
		file.WriteString(fmt.Sprintf(vertex, x, y, z))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y, z))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y, z + cubeSize))
		file.WriteString(fmt.Sprintf(endLoop))
		file.WriteString(fmt.Sprintf(endFacet))

		file.WriteString(fmt.Sprintf(facetNormal, 0.0,-1.0,0.0))
		file.WriteString(fmt.Sprintf(outerLoop))
		file.WriteString(fmt.Sprintf(vertex, x, y, z))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y, z + cubeSize))
		file.WriteString(fmt.Sprintf(vertex, x, y, z + cubeSize))
		file.WriteString(fmt.Sprintf(endLoop))
		file.WriteString(fmt.Sprintf(endFacet))

		file.WriteString(fmt.Sprintf(facetNormal, 0.0,0.0,1.0))
		file.WriteString(fmt.Sprintf(outerLoop))
		file.WriteString(fmt.Sprintf(vertex, x, y, z + cubeSize))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y, z + cubeSize))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y + cubeSize, z + cubeSize))
		file.WriteString(fmt.Sprintf(endLoop))
		file.WriteString(fmt.Sprintf(endFacet))

		file.WriteString(fmt.Sprintf(facetNormal, 0.0,0.0,1.0))
		file.WriteString(fmt.Sprintf(outerLoop))
		file.WriteString(fmt.Sprintf(vertex, x, y, z + cubeSize))
		file.WriteString(fmt.Sprintf(vertex, x + cubeSize, y + cubeSize, z + cubeSize))
		file.WriteString(fmt.Sprintf(vertex, x, y + cubeSize, z))
		file.WriteString(fmt.Sprintf(endLoop))
		file.WriteString(fmt.Sprintf(endFacet))
}