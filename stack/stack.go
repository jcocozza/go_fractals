package stack

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/jcocozza/go_fractals/utils"
	IFS "github.com/jcocozza/go_fractals/IteratedFunctionSystems"
	viz "github.com/jcocozza/go_fractals/visualizer"
)

func CreateFractalStack(ifspath string, numStacks int, thickness float32, outputPath string) {
	transformList, dim := IFS.ParseIFS(ifspath)
	newIfs := IFS.NewIteratedFunctionSystem(transformList, 3, 1000, dim)

	base := newIfs.InitialPoints
	baseImg := viz.NewFractalImage(1000,1000,"base.png", base)

	var tempPoints = newIfs.InitialPoints
	var tempFractalImg *viz.FractalImage
	var fractalImageList []*viz.FractalImage

	fractalImageList = append(fractalImageList, baseImg)
	for i := 0; i < numStacks; i++ {
		tempPoints = newIfs.RunDeterministicStepWise(tempPoints)
		tempFractalImg = viz.NewFractalImage(1000,1000,fmt.Sprintf("stack-%d",i), tempPoints)
		fractalImageList = append(fractalImageList, tempFractalImg)
	}

	var mesh [][]float32
	var tmpMesh [][]float32

	var start float32 = 0
	for _, fractal := range fractalImageList {
		tmpMesh = fractal.ExtrudeImage(start,start+thickness)
		mesh = append(mesh, tmpMesh...)
		start = start + thickness
	}

	saveSTL(mesh, outputPath)

}

func crossProduct(v1, v2 []float32) []float32 {
	return []float32{
		v1[1]*v2[2] - v1[2]*v2[1],
		v1[2]*v2[0] - v1[0]*v2[2],
		v1[0]*v2[1] - v1[1]*v2[0],
	}
}
func saveSTL(mesh [][]float32, filePath string) {
	// Create an STL file with the given mesh
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating STL file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write the binary STL header (80 bytes)
	header := make([]byte, 80)
	file.Write(header)

	// Write the number of facets (triangles) as a 32-bit unsigned integer
	numFacets := uint32(len(mesh) / 3)
	binary.Write(file, binary.LittleEndian, numFacets)

	// Write each facet to the STL file
	for i := 0; i < len(mesh); i += 3 {
		// Calculate the normal vector for each facet
		v1 := []float32{mesh[i+1][0] - mesh[i][0], mesh[i+1][1] - mesh[i][1], mesh[i+1][2] - mesh[i][2]}
		v2 := []float32{mesh[i+2][0] - mesh[i][0], mesh[i+2][1] - mesh[i][1], mesh[i+2][2] - mesh[i][2]}
		normal := crossProduct(v1, v2)

		// Write the normal vector components to the STL file
		for _, component := range normal {
			binary.Write(file, binary.LittleEndian, component)
		}

		// Write the vertices of the facet to the STL file
		for j := 0; j < 3; j++ {
			for _, component := range mesh[i+j] {
				binary.Write(file, binary.LittleEndian, component)
			}
		}

		// Write the attribute byte count (set to zero for simplicity)
		binary.Write(file, binary.LittleEndian, uint16(0))
		utils.ProgressBar(i+1, len(mesh))
	}
}

