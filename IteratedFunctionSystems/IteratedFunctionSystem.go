package IteratedFunctionSystems

import (
	"fmt"
	"math/rand"
	"math"
	"time"

	"github.com/jcocozza/go_fractals/utils"
	"gonum.org/v1/gonum/mat"
)


//An IteratedFunctionSystem is a list of Transformations
type IteratedFunctionSystem struct {
	TransformationList []Transformation
	NumIterations int
	InitialPoints [][]float64
}

//NewIteratedFunctionSystem creates an iterated function system. It defaults to 1 initial point
func NewIteratedFunctionSystem(transformationList []Transformation, numIterations int, numPoints int, numDims int) (*IteratedFunctionSystem) {
	return &IteratedFunctionSystem{
		TransformationList: transformationList,
		NumIterations: numIterations,
		InitialPoints: _GenerateInitialPoints(numPoints, numDims),
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
	//fmt.Println("Created Transformation Points:", pointTransformationList)
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
		if len(ifs.TransformationList) == 1{
			pointsList = append(pointsList, tempPointList...) // allow the deterministic algorithm to work with 1 transformation
		} else {
			pointsList = tempPointList
		}
		//fmt.Println("Point Space:", pointsList)
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Println("Total number of points:", len(pointsList))
	//fmt.Println("Total number of Unique points:", utils.CalculateUniqueListElements(pointsList))
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

// Calculate the probabilities to use a given matrix in the Probabilistic method
/*

prob_weights = DET(M_1) / DET(M_1) + DET(M_2) + ...

Returns a list of probabilities ordered by the order of transformations passed in
*/
func (ifs *IteratedFunctionSystem) CalculateProbabilities() []float64 {

	determinantSum := 0.0
	var probabilities []float64
	for _, transform := range ifs.TransformationList {
		determinantSum = determinantSum + math.Abs(mat.Det(&transform.Similarity))
	}

	for _, transform := range ifs.TransformationList {
		probability := math.Abs(mat.Det(&transform.Similarity)) / determinantSum
		probabilities = append(probabilities, probability)
	}

	return probabilities
}

// Calculate probabilities -- possibly a way to handle to determinant of 0 problem
func (ifs *IteratedFunctionSystem) CalculateProbabilitiesTEST() []float64 {
    determinantSum := 0.0
    epsilon := 1e-6 // Small epsilon value to ensure nonzero probabilities

    var probabilities []float64
    for _, transform := range ifs.TransformationList {
        determinant := math.Abs(mat.Det(&transform.Similarity))
        determinantSum += determinant
    }

    for _, transform := range ifs.TransformationList {
        determinant := math.Abs(mat.Det(&transform.Similarity))
        probability := (determinant + epsilon) / (determinantSum + epsilon*float64(len(ifs.TransformationList)))
        probabilities = append(probabilities, probability)
    }

    return probabilities
}

//runs the probabilistic algorithm for an Iterated Function System
func (ifs *IteratedFunctionSystem) RunProbabilistic(probabilities []float64) [][]float64 {
	startTime := time.Now()
	pointsList := ifs.InitialPoints
	mostRecentPoints := ifs.InitialPoints
	fmt.Println("probabilities:", probabilities)
	for i := 0; i < ifs.NumIterations; i++ {
		for j := 0; j < len(mostRecentPoints); j++ {
			var newPoint []float64
			idx := utils.PickRandom(probabilities) //pick an index based on the probabilities
			newPoint, _ = ifs.TransformationList[idx].Transform(mostRecentPoints[j]) // get the transform at that index and apply to a point
			pointsList = append(pointsList, newPoint) // keep track of all the points; we will eventually graph this
			mostRecentPoints[j] = newPoint // we only need to apply the transformations to the next point(s) produced.
		}
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Println("Total number of points:", len(pointsList))
	//fmt.Println("Total number of Unique points:", utils.CalculateUniqueListElements(pointsList))
	fmt.Printf("Elapsed time for Probabilistic algorithm: %v\n", elapsedTime)
	return pointsList
}

func (ifs *IteratedFunctionSystem) RunProbabilisticStepWise(stepPoints [][]float64) [][]float64 {
	mostRecentPoints := stepPoints
	probabilities := ifs.CalculateProbabilities() // TODO remove this so it doesn't need to be calculated every time

	for j := 0; j < len(mostRecentPoints); j++ {
		var newPoint []float64
		idx := utils.PickRandom(probabilities) //pick an index based on the probabilities
		newPoint, _ = ifs.TransformationList[idx].Transform(mostRecentPoints[j]) // get the transform at that index and apply to a point
		mostRecentPoints[j] = newPoint // we only need to apply the transformations to the next point(s) produced.
	}
	return mostRecentPoints // we only want the most recent set of points
}