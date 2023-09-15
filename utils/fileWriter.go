package utils

import (
    "os"
	"fmt"
)

//write the contents of a variable to a file
func Writer(fileName string,data interface{}) {
    // Create or open the file for writing
    file, err := os.Create(fileName)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // Variable with content to write
	content := fmt.Sprintf("%v", data)
    // Write the content to the file
    _, err = file.Write([]byte(content))
    if err != nil {
        panic(err)
    }
}
