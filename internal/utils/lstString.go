package utils

import "fmt"

// convert []float64{a,b,c} to a string "a,b,c"
func ListToString(lst []float64) string {
	str := ""

	for i, elm := range lst {
		if i + 1 != len(lst) {
			tstr := fmt.Sprintf("%f", elm) + ","
			str += tstr
		} else {
			tstr := fmt.Sprintf("%f", elm)
			str += tstr
		}
	}
	return str

}

// convert []string{"a","b","c"} to a string "a,b,c"
func StrListToString(lst []string) string {
	str := ""

	for i, elm := range lst {
		if i+1 != len(lst) {
			str += elm + " "
		} else {
			str += elm
		}
	}
	return str
}