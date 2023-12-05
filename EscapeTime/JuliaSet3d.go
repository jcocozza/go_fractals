package EscapeTime

import (
	"encoding/binary"
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

	HeaderSize  = 80
)

// function to get alpha value of a pixel, handling out-of-bounds gracefully
func getAlpha(img image.Image, x, y int) uint32 {
	bounds := img.Bounds()
	if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
		pixel := img.At(x, y)
		_, _, _, alpha := pixel.RGBA()
		return alpha
	}
	return 0
}

// if any adjacent points are transparent, then the non-transparent pixel is on the boundary
func isTransition(img image.Image, x,y int) bool {
	alphaLeft := getAlpha(img, x-1, y)
	alphaRight := getAlpha(img, x+1, y)
	alphaUp := getAlpha(img, x, y-1)
	alphaDown := getAlpha(img, x, y+1)

	return alphaLeft == 0 || alphaRight == 0 || alphaUp == 0 || alphaDown == 0
}

// check if the pixel in the current image will be "covered" by pixels in other images
// the specific use case is in an ordered list of images
// x,y will be apoint in the current image, and images will be a list of future images
func isCovered(x,y int, images []image.Image) bool {
	for _, img := range images {
		pixel := img.At(x, y)
			_, _, _, alpha := pixel.RGBA()
			if alpha > 0 {
				return true
			}
	}
	return false
}

// for each non-transparent pixel in the image, draw a thickend, filled julia set
func DrawJuliaSet3DFilled(img image.Image, stlFile *os.File, shift float64) {
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


// for each non-transparent pixel in the image, draw a thickend, outline of the set
// images is the list of images that will be stacked on top of the passed img.
func DrawJuliaSet3DEmpty(img image.Image, images []image.Image, stlFile *os.File, shift float64) {
	// Loop through pixels and generate cuboids for non-transparent pixels
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			_, _, _, alpha := pixel.RGBA()

			if alpha > 0 {
				if !isCovered(x,y, images) {
					writeCube(stlFile, float64(x), float64(y), shift)
				} else if isTransition(img, x,y) {
					writeCube(stlFile, float64(x), float64(y), shift)
				}
			}
		}
	}
}


func DrawJuliaSet3DBinary(img image.Image, stlFile *os.File, shift float64) {
	// Loop through pixels and generate cuboids for non-transparent pixels
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			_, _, _, alpha := pixel.RGBA()

			if alpha > 0 {
				WriteCubeBinary(stlFile, float32(x), float32(y), float32(shift))
			}
		}
	}
}

func writeNormalBinary(file *os.File, normal [3]float32) {
	binary.Write(file, binary.LittleEndian, &normal[0])
	binary.Write(file, binary.LittleEndian, &normal[1])
	binary.Write(file, binary.LittleEndian, &normal[2])
}

func writeFacetBinary(file *os.File, vertices [3][3]float32) {
	for _, vertex := range vertices {
		for _, value := range vertex {
			binary.Write(file, binary.LittleEndian, &value)
		}
	}
	file.Write([]byte{0, 0})
}

// draw a cube at an (x,y,z) coordinate in binary
func WriteCubeBinary(file *os.File, x,y,z float32) {
	var normal [3]float32
	var vertices [3][3]float32

	normal = [3]float32{0.0, 0.0, -1.0}
	vertices = [3][3]float32{
		{x, y, z},
		{x + cubeSize, y + cubeSize, z},
		{x + cubeSize, y, z},
	}
	writeNormalBinary(file, normal)
	writeFacetBinary(file, vertices)

	normal = [3]float32{0.0, 0.0, -1.0}
	vertices = [3][3]float32{
		{x, y, z},
		{x, y + cubeSize, z},
		{x + cubeSize, y + cubeSize, z},
	}
	writeNormalBinary(file, normal)
	writeFacetBinary(file, vertices)

	normal = [3]float32{-1.0,0.0,0.0}
	vertices = [3][3]float32{
		{x, y, z},
		{x, y + cubeSize, z + cubeSize},
		{x, y + cubeSize, z},
	}
	writeNormalBinary(file, normal)
	writeFacetBinary(file, vertices)

	normal = [3]float32{-1.0,0.0,0.0}
	vertices = [3][3]float32{
		{x, y, z},
		{x, y, z + cubeSize},
		{x, y + cubeSize, z + cubeSize},
	}
	writeNormalBinary(file, normal)
	writeFacetBinary(file, vertices)

	normal = [3]float32{0.0,1.0,0.0}
	vertices = [3][3]float32{
		{x, y + cubeSize, z},
		{x + cubeSize, y + cubeSize, z + cubeSize},
		{x + cubeSize, y + cubeSize, z},
	}
	writeNormalBinary(file, normal)
	writeFacetBinary(file, vertices)

	normal = [3]float32{0.0,1.0,0.0}
	vertices = [3][3]float32{
		{x, y + cubeSize, z},
		{x, y + cubeSize, z + cubeSize},
		{x + cubeSize, y + cubeSize, z + cubeSize},
	}
	writeNormalBinary(file, normal)
	writeFacetBinary(file, vertices)

	normal = [3]float32{1.0,0.0,0.0}
	vertices = [3][3]float32{
		{x + cubeSize, y, z},
		{x + cubeSize, y + cubeSize, z},
		{x + cubeSize, y + cubeSize, z + cubeSize},
	}
	writeNormalBinary(file, normal)
	writeFacetBinary(file, vertices)

	normal = [3]float32{1.0,0.0,0.0}
	vertices = [3][3]float32{
		{x + cubeSize, y, z},
		{x + cubeSize, y + cubeSize, z + cubeSize},
		{x + cubeSize, y, z + cubeSize},
	}
	writeNormalBinary(file, normal)
	writeFacetBinary(file, vertices)

	normal = [3]float32{0.0,-1.0,0.0}
	vertices = [3][3]float32{
		{x, y, z},
		{x + cubeSize, y, z},
		{x + cubeSize, y, z + cubeSize},
	}
	writeNormalBinary(file, normal)
	writeFacetBinary(file, vertices)

	normal = [3]float32{0.0,-1.0,0.0}
	vertices = [3][3]float32{
		{x, y, z},
		{x + cubeSize, y, z + cubeSize},
		{x, y, z + cubeSize},
	}
	writeNormalBinary(file, normal)
	writeFacetBinary(file, vertices)

	normal = [3]float32{0.0,0.0,1.0}
	vertices = [3][3]float32{
		{x, y, z + cubeSize},
		{x + cubeSize, y, z + cubeSize},
		{x + cubeSize, y + cubeSize, z + cubeSize},
	}
	writeNormalBinary(file, normal)
	writeFacetBinary(file, vertices)

	normal = [3]float32{0.0,0.0,1.0}
	vertices = [3][3]float32{
		{x, y, z + cubeSize},
		{x + cubeSize, y + cubeSize, z + cubeSize},
		{x, y + cubeSize, z},
	}
	writeNormalBinary(file, normal)
	writeFacetBinary(file, vertices)

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