package utils

import (
	"fmt"
)

// Determine the number of unique subLists in a list of lists
func CalculateUniqueListElements(lst [][]float64) int{
	uniqueElements := make(map[string]bool)

    // Iterate through the list and add serialized subList to the map
    for _, subList := range lst {
        subListStr := fmt.Sprintf("%v", subList)

                // Check if it's unique, and if so, add it to the map
        if _, exists := uniqueElements[subListStr]; !exists {
            uniqueElements[subListStr] = true
        }
    }
	return len(uniqueElements) // the number of distinct subLists
}