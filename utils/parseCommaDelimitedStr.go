package utils

import (
	"strconv"
	"strings"
	"fmt"
)

//Return a list of float64 from a string list representation:
//"a,b,c,d" -> [a,b,c,d]
func ParseCommaDelimitedStr(cds string) []float64 {
	elements := strings.Split(cds, ",")

	var data []float64

	for _, element := range elements {
		num, err := strconv.ParseFloat(element, 64)
		if err != nil {
			fmt.Println("Error parsing element:", err)
			continue // TODO this currently just skips this element if there's an error
		}
		data = append(data, num)
	}
	return data
}