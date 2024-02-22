package controllers

import "fmt"

type ComplexFromTS struct {
	Real float64 `json:"real"`
	Imaginary float64 `json:"imaginary"`
}
func (cts *ComplexFromTS) ToComplex() complex128 {
	return complex(cts.Real, cts.Imaginary)
}
func (cts *ComplexFromTS) ToString() string {
	if cts.Imaginary < 0 {
		return fmt.Sprint(cts.Real) + fmt.Sprint(cts.Imaginary) + "i"
	}
	return fmt.Sprint(cts.Real) + "+" + fmt.Sprint(cts.Imaginary) + "i"
}