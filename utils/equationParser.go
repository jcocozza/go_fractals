package utils

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"plugin"
	"strings"
)

/*
This allow a user to attach arbritary functions via a string
*/

const (
	template = `package main

%s

var ParsedTransformation = func(%s) complex128 {
	return %s
}
`

	importCmplx = `import "math/cmplx"`
	importMath = `import "math"`

	twoParam = `z,c complex128`
	oneParam = `z complex128`
)

type TwoParamEquation func(complex128,complex128) complex128
type OneParamEquation func(complex128) complex128

type ParsedEquation struct {
    Equation interface{}
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

func getSym(funcString string) plugin.Symbol {
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
	return sym
}

// handle 1 parameter functions
func CreateOneParamEquation(eqnInput string) OneParamEquation {
	var sym plugin.Symbol
    newEqnString := template

	if strings.Contains(eqnInput, "cmplx") {
        slog.Info("Equation contains 'cmplx'")
        newEqnString = fmt.Sprintf(newEqnString, importCmplx, oneParam, eqnInput)
    } else if strings.Contains(eqnInput, "math") {
		slog.Info("Equation contains 'math'")
		newEqnString = fmt.Sprintf(newEqnString, importMath, oneParam, eqnInput)
	} else {
        newEqnString = fmt.Sprintf(newEqnString, "", oneParam, eqnInput) // no imports
    }
    funcString := newEqnString
    sym = getSym(funcString)
    parsedTransformationFunc, ok := sym.(*func(complex128) complex128)

	if !ok {
		fmt.Println("Unexpected type for symbol")
		return nil
	}
	return *parsedTransformationFunc
}

// handle 2 parameter function
func CreateTwoParamEquation(eqnInput string) TwoParamEquation {
	var sym plugin.Symbol
    newEqnString := template

	if strings.Contains(eqnInput, "cmplx") {
        slog.Info("Equation contains 'cmplx'")
        newEqnString = fmt.Sprintf(newEqnString, importCmplx, twoParam, eqnInput)
    } else if strings.Contains(eqnInput, "math") {
		slog.Info("Equation contains 'math'")
		newEqnString = fmt.Sprintf(newEqnString, importMath, twoParam, eqnInput)
	} else {
        newEqnString = fmt.Sprintf(newEqnString, "", twoParam, eqnInput) // no imports
    }

	slog.Info("Equation has 2 variables")
	funcString := newEqnString
	fmt.Println(funcString)
	sym = getSym(funcString)
	parsedTransformationFunc, ok := sym.(*func(complex128,complex128) complex128)

	if !ok {
		fmt.Println("Unexpected type for symbol")
		return nil
	}
	return *parsedTransformationFunc
}



/*
// Parse equations with 2 parameters
func ParseEquation2(eqnInput string) func(complex128, complex128) complex128 {
	funcString := fmt.Sprintf(`
package main

var ParsedTransformation = func (z,c complex128) complex128 {
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
	parsedTransformationFunc, ok := sym.(*func(complex128,complex128) complex128)
	if !ok {
		fmt.Println("Unexpected type for symbol")
		return nil
	}
	return *parsedTransformationFunc
}

// parse equations with 1 parameter
func ParseEquation(eqnInput string) func(complex128) complex128 {
	//eqnInput := "1/(c*c + .72i)" // Replace this with your equation

	// Dynamically create the function string
	funcString := fmt.Sprintf(`
package main

var ParsedTransformation = func(z complex128) complex128 {
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

	return *parsedTransformationFunc
}
*/

