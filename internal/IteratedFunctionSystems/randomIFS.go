package IteratedFunctionSystems

import (
	"log/slog"
	"math/rand"

	"github.com/jcocozza/go_fractals/internal/utils"
	"gonum.org/v1/gonum/mat"
)

func generateRandomMatrix(rows int, cols int) mat.Dense {
	var matrixData []float64
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			matrixData = append(matrixData, rand.Float64()*2 - 1)
		}
	}
	slog.Info("New Matrix generated: " + utils.ListToString(matrixData))
	return *mat.NewDense(rows, cols, matrixData)

}

func generateRandomTransformation(dimension int) Transformation {
	similarity := generateRandomMatrix(dimension, dimension)
	shift := generateRandomMatrix(dimension, 1)
	randomTransform, _ := NewTransformation(similarity, shift)
	return *randomTransform
}

func GenerateRandomIFS(dimension int, numTransformations int) IteratedFunctionSystem {
	var transformationList []Transformation

	for i := 0; i < numTransformations; i ++ {
		tempTransform := generateRandomTransformation(dimension)
		transformationList = append(transformationList, tempTransform)
	}

	return *NewIteratedFunctionSystem(transformationList, 1000000, 1, dimension)
}
