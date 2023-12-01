package utils

import (
	"fmt"
	"os"
	"os/exec"
	"plugin"
)

/*
var juliaEquation = `var ParsedTransformation = func(c complex128) complex128 {return %s}`
var juliaEvolutionEquation = `var ParsedTransformation = func(z complex128, shiftParam complex128) complex128 {return %s}`
var mandelbrotEquation = `var ParsedTransformation = func(z, c complex128) complex128 {return %s}`
*/
func ParseEquation(eqnInput string) func(complex128) complex128 {
	//eqnInput := "1/(c*c + .72i)" // Replace this with your equation

	// Dynamically create the function string
	funcString := fmt.Sprintf(`
package main

var ParsedTransformation = func(c complex128) complex128 {
	return %s
}`, eqnInput)

	// Write the function string to a temporary file
	fileName := "dynamic_function.go"
	err := writeToFile(fileName, funcString)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return nil
	}
	defer os.Remove(fileName)

	// Compile the Go file to a shared object (.so) file
	soFileName := "dynamic_function.so"
	err = compileToSharedObject(fileName, soFileName)
	if err != nil {
		fmt.Println("Error compiling to shared object:", err)
		return nil
	}
	defer os.Remove(soFileName)

	// Load the plugin
	p, err := plugin.Open(soFileName)
	if err != nil {
		fmt.Println("Error opening plugin:", err)
		return nil
	}

	// Look up the symbol (function) from the plugin
	sym, err := p.Lookup("ParsedTransformation")
	if err != nil {
		fmt.Println("Error looking up symbol:", err)
		return nil
	}

	// Assert the symbol to the expected type
	parsedTransformationFunc, ok := sym.(*func(complex128) complex128)
	if !ok {
		fmt.Println("Unexpected type for symbol")
		return nil
	}

	// Now you can use the function
	//result := (*parsedTransformationFunc)(2 + 3i)
	return *parsedTransformationFunc
}

func writeToFile(fileName, content string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func compileToSharedObject(inputFile, outputFile string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outputFile, inputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}