package primitives

import(
	"fmt"
	"gonum.org/v1/gonum/mat"
)


type vectorSpace struct {
	Dimension int
	UnitBasis mat.Dense
}

/*
Create a vector space
*/
func newVectorSpace(dimension int) (*vectorSpace) {
	var unitBasis []float64
	for i :=0; i<dimension; i++ {
		vec := make([]float64, dimension)
		vec[i] = 1
		unitBasis = append(unitBasis, vec...)
	}
	return &vectorSpace{
		Dimension: dimension,
		UnitBasis: *mat.NewDense(dimension,dimension,unitBasis),
	}
}

func Test() {
	var result mat.Dense
	var result2 mat.Dense
	vs := newVectorSpace(3)
	vs2 := newVectorSpace(3)
	result.Add(&vs.UnitBasis,&vs2.UnitBasis)
	result2.Mul(&result,&result)
	fmt.Println(result2)
}
