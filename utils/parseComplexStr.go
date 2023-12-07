package utils

import (
	"regexp"
	"strconv"
)

// parse the "a+bi" string into complex(a,b)
func ParseComplexString(input string) complex128 {
	// Define a regular expression to extract real and imaginary parts
	complexPattern := regexp.MustCompile(`^([-+]?\d*\.?\d+)([-+]\d*\.?\d*)i$`)

	// Find the matched groups in the input string
	matches := complexPattern.FindStringSubmatch(input)
	if len(matches) != 3 {
		panic("invalid complex number format")
	}

	// Extract real and imaginary parts from the matched groups
	realPart, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		//"error parsing real part: %v"
		panic(err)
	}

	imagPart, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		panic (err)
		//"error parsing imaginary part: %v"
	}

	// Create a complex number using the parsed real and imaginary parts
	result := complex(realPart, imagPart)
	return result
}