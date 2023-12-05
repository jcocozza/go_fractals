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