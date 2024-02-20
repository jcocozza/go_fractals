package utils
import (
	"math/rand"
)

//Return the index of the element to use in the list
func PickRandom(probabilities []float64) int {
	r := rand.Float64()
	cumulative := 0.0

	for i, probability := range probabilities {
		cumulative = cumulative + probability

		if r < cumulative {
			return i
		}
	}
	panic("YIKES! failed to pick random number")
}