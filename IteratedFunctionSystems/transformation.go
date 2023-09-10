package IteratedFunctionSystems

import (
	"fmt"
    "gonum.org/v1/gonum/mat"
)

//A Transformation consisting of a similarity and a shift
type Transformation struct {
	Similarity mat.Dense
	Shift mat.Dense
}

//NewTransformation creates a new transformation
func NewTransformation(similarity mat.Dense, shift mat.Dense) (*Transformation, error) {
	numRowsSI, numColsSI := similarity.Dims()
	numRowsSH, numColsSH := shift.Dims()

	// the similarity needs to be square since we are operating on basis vectors
	// the shift needs to be the same length as the number of columns in the similarity
	// A shift can only have 1 column
	if numRowsSI != numColsSI || numRowsSH != numColsSI || numColsSH != 1 {
		return &Transformation{}, fmt.Errorf("cannot do transform on improperly shaped data")
	}

	return &Transformation{
		Similarity: similarity,
		Shift: shift,
	}, nil
}

//Transform performs the transformation: similarity * point + shift
func (t *Transformation) Transform(point []float64) ([]float64, error) {
	// if we don't put in the same number of points as there are columns then we can't do the calculation
	_, numCols := t.Similarity.Dims()
	if len(point) != numCols {
		return nil, fmt.Errorf("cannot do transform on improperly shaped data")
	}

	m := mat.NewDense(len(point), 1, point) // the result will only have 1 column, since the point only has 1 column
	var result mat.Dense // to calculate the result, multiply the point by the similarity, then add the shift
	result.Mul(&t.Similarity, m)
	result.Add(&result, &t.Shift)
	return result.RawMatrix().Data, nil
}
