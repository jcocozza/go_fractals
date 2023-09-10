package IteratedFunctionSystems

import (
	"fmt"
	"time"
	"math/rand"
)

//An IteratedFunctionSystem is a list of Transformations
type IteratedFunctionSystem struct {
	TransformationList []Transformation
	NumIterations int
	InitialPoints [][]float64
}

//NewIteratedFunctionSystem creates an iterated function system. It defaults to 1 initial point
func NewIteratedFunctionSystem(transformationList []Transformation, numIterations int, numDims int) (*IteratedFunctionSystem) {
	return &IteratedFunctionSystem{
		TransformationList: transformationList,
		NumIterations: numIterations,
		InitialPoints: _GenerateInitialPoints(10, numDims),
	}
}

/*
_GenerateInitialPoints creates the initial points for an Iterated Function system. By default this is a random coordinate tuple
where each entry is between 0 and 1.
*/
func _GenerateInitialPoints(numPoints int, numDims int) [][]float64 {
	initialPoints := make([][]float64, numPoints)

	// Generate random initial points within a reasonable range
	for i := 0; i < numPoints; i++ {
		point := make([]float64, numDims)
		for j := 0; j < numDims; j++ {
			point[j] = rand.Float64() // Random value between 0 and 1 for each dimension
		}
		initialPoints[i] = point
	}

	return initialPoints
}

/*
Transform performs an Iterated Function System Transformation.
It takes a single point and returns the list of points corresponding to each transformation in the system.
*/
func (ifs *IteratedFunctionSystem) Transform(point []float64) [][]float64 {
	pointTransformationList := make([][]float64, len(ifs.TransformationList))

	for i := 0; i < len(ifs.TransformationList); i++ {
		tempPoint, _ := ifs.TransformationList[i].Transform(point)
		pointTransformationList[i] = tempPoint
	}
	fmt.Println("Created Transformation Points:", pointTransformationList)
	return pointTransformationList
}

//RunDeterministic runs the deterministic algorithm for an Iterated Function System
func (ifs *IteratedFunctionSystem) RunDeterministic() [][]float64 {
	startTime := time.Now()
	pointsList := ifs.InitialPoints

	for i := 0; i < ifs.NumIterations; i++ {
		var tempPointList [][]float64
		for j := 0; j < len(pointsList); j++ {
			newPoints := ifs.Transform(pointsList[j])
			tempPointList = append(tempPointList, newPoints...)
		}
		pointsList = tempPointList
		//fmt.Println("Point Space:", pointsList)
	}
	endTime := time.Now()

	elapsedTime := endTime.Sub(startTime)
	fmt.Println("Total number of points:", len(pointsList))
	fmt.Printf("Elapsed time for Deterministic algorithm: %v\n", elapsedTime)
	return pointsList
}

//RunDeterministicStepWise will run a single iteration of the RunDeterministic loop
func (ifs *IteratedFunctionSystem) RunDeterministicStepWise(stepPoints [][]float64) [][]float64 {
	pointsList := stepPoints
	var finalPointList [][]float64
	for i := 0; i < len(pointsList); i++ {
		newPoints := ifs.Transform(pointsList[i])
		finalPointList = append(finalPointList, newPoints...)
	}
	return finalPointList
}