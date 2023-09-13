package utils

import (
	"fmt"
	"encoding/json"
)

// Determine the number of unique subLists in a list of lists
func CalculateUniqueListElements(lst [][]float64) int{
	uniqueElements := make(map[string]bool)

    // Iterate through the list and add serialized subList to the map
    for _, subList := range lst {
        subListStr, err := json.Marshal(subList)
        if err != nil {
            fmt.Println("Error:", err)
            panic("")
        }
        uniqueElements[string(subListStr)] = true
    }
	return len(uniqueElements) // the number of distinct subLists
}
