package images

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"log/slog"
	"os"

	"github.com/jcocozza/go_fractals/internal/utils"
)

const (
	cubeSize    = 1.0
	facetNormal = "  facet normal %f %f %f\n"
	outerLoop   = "    outer loop\n"
	vertex      = "      vertex %f %f %f\n"
	endLoop     = "    endloop\n"
	endFacet    = "  endfacet\n"

	headerSize  = 80
)

type STLWriteMethod func(*os.File,float64,float64,float64)
type ExtrusionMethod func(image.Image,[]image.Image, *os.File, float64, STLWriteMethod) int

// master func for create stl's
func STLControlFlow(writeBinary, solid bool, imgList []image.Image, filePath string) {
	if writeBinary {
		if solid {
			createSTLBinary(imgList, filePath, extrudeImageSolid)
		} else {
			createSTLBinary(imgList, filePath, extrudeImageHollow)
		}
	} else {
		if solid {
			createSTL(imgList, filePath, extrudeImageSolid)
		} else {
			createSTL(imgList, filePath, extrudeImageHollow)
		}
	}
}

func createSTL(imgList []image.Image, filePath string, extrustionMethod ExtrusionMethod) {
	stlFile, err := os.Create(filePath)
	if err != nil {
		slog.Error("Error creating STL file:", err)
		return
		}
	defer stlFile.Close()

	shift := 0.0
	stlFile.WriteString("solid GeneratedModel\n")
	for i, img := range imgList {
		extrustionMethod(img, imgList[i+1:], stlFile, shift, writeCube)
		shift += 1

		utils.ProgressBar(i,len(imgList))
	}
	stlFile.WriteString("endsolid GeneratedModel\n")
}

func createSTLBinary(imgList []image.Image, filePath string, extrustionMethod ExtrusionMethod) {
	stlFile, err := os.Create(filePath)
	if err != nil {
		slog.Error("Error creating STL file:", err)
		return
		}
	defer stlFile.Close()

	header := make([]byte, headerSize)
	stlFile.Write(header)
	// Record the position after writing the header
	facetCountPos, err := stlFile.Seek(0, io.SeekCurrent)
	if err != nil {
		slog.Error("Error seeking file:", err)
		return
	}
	// Write a temporary placeholder for the facet count
	tmpFacetCount := make([]byte, 4)
	stlFile.Write(tmpFacetCount)

	shift := 0.0
	totalFacets := 0
	for i, img := range imgList {
		totalFacets += extrustionMethod(img, imgList[i+1:], stlFile, shift, writeCubeBinary)
		shift += 1
		utils.ProgressBar(i,len(imgList))
	}


	// Move back to the position to write the actual facet count
	_, err = stlFile.Seek(facetCountPos, io.SeekStart)
	if err != nil {
		slog.Error("Error seeking file:", err)
		return
	}

	// Write the actual facet count
	binary.LittleEndian.PutUint32(tmpFacetCount, uint32(totalFacets))
	stlFile.Write(tmpFacetCount)


	// Move back to the end of the file
	_, err = stlFile.Seek(0, io.SeekEnd)
	if err != nil {
		slog.Error("Error seeking file:", err)
		return
	}
}

// for each non-transparent pixel in the image, draw a thickend set
func extrudeImageSolid(img image.Image, images []image.Image, stlFile *os.File, shift float64, method STLWriteMethod) int {
	bounds := img.Bounds()
	var totalFacets int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			_, _, _, alpha := pixel.RGBA()

			if alpha > 0 {
				method(stlFile, float64(x), float64(y), shift)
				totalFacets += 12 //every square has 12 facets
			}
		}
	}
	return totalFacets
}

// for each non-transparent pixel in the image, draw a thickend, outline of the set
// images is the list of images that will be stacked on top of the passed img.
func extrudeImageHollow(img image.Image, images []image.Image, stlFile *os.File, shift float64, method STLWriteMethod) int {
	bounds := img.Bounds()
	var totalFacets int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			_, _, _, alpha := pixel.RGBA()

			if alpha > 0 {
				if !IsCovered(x,y, images) {
					method(stlFile, float64(x), float64(y), shift)
					totalFacets += 12
				} else if IsTransition(img, x,y) {
					method(stlFile, float64(x), float64(y), shift)
					totalFacets += 12
				}
			}
		}
	}
	return totalFacets
}

// helper for writeCubeBinary()
func writeNormalBinary(file *os.File, normal [3]float32) {
	binary.Write(file, binary.LittleEndian, &normal[0])
	binary.Write(file, binary.LittleEndian, &normal[1])
	binary.Write(file, binary.LittleEndian, &normal[2])
}

// helper for writeCubeBinary()
func writeFacetBinary(file *os.File, vertices [3][3]float32) {
	for _, vertex := range vertices {
		for _, value := range vertex {
			binary.Write(file, binary.LittleEndian, &value)
		}
	}
	file.Write([]byte{0, 0})
}

// draw a cube at an (x,y,z) coordinate in binary
func writeCubeBinary(file *os.File, x64,y64,z64 float64) {

	x := float32(x64)
	y := float32(y64)
	z := float32(z64)

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