package IteratedFunctionSystems

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/jcocozza/go_fractals/internal/utils"
	"gonum.org/v1/gonum/mat"
)

func ParseIFS(filePath string) ([]Transformation, int) {
	slog.Debug("Reading from file: " + filePath)
	rowMatches := utils.Reader(filePath, `\[(.*?)\]`)

	var transformationList []Transformation
	var dimension int
	for _, rowMatch := range rowMatches {
		dimSimilarity := utils.ParseCommaDelimitedStr(strings.Trim(rowMatch[0], "[]"))
		similarity := utils.ParseCommaDelimitedStr(strings.Trim(rowMatch[1], "[]"))

		dimShift := utils.ParseCommaDelimitedStr(strings.Trim(rowMatch[2], "[]"))
		shift := utils.ParseCommaDelimitedStr(strings.Trim(rowMatch[3], "[]"))

		dimension = int(dimShift[0])

		similarityMatrix := mat.NewDense(int(dimSimilarity[0]),int(dimSimilarity[1]), similarity)
		shiftMatrix := mat.NewDense(int(dimShift[0]),int(dimShift[1]), shift)
		transform, err := NewTransformation(*similarityMatrix, *shiftMatrix)

		if err != nil {
			panic(fmt.Sprintf("Unable to create transform: %s", err))
		}

		transformationList = append(transformationList, *transform)
	}
	return transformationList, dimension
}