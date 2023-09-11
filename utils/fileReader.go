package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

/*
Return a list of regex matches line by line - empty list for no matches
Each element in the list is a list of matches by line
*/
func Reader(path string, pattern string) [][]string {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error Opening File", err)
		return [][]string{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	patternRegEx := regexp.MustCompile(pattern)

	//check each line for the pattern
	var matches [][]string
	for scanner.Scan() {
		line := scanner.Text()

		matches = append(matches, patternRegEx.FindAllString(line, -1))
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return [][]string{}
	}
	return matches
}