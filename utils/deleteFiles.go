package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func DeleteFiles(dirPath string, filePattern string) {
    // Use filepath.Glob to get a list of matching file paths
    matchingFiles, err := filepath.Glob(filepath.Join(dirPath, filePattern))
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Loop through the matching files and delete them
    for _, filePath := range matchingFiles {
        err := os.Remove(filePath)
        if err != nil {
            fmt.Printf("Error deleting %s: %v\n", filePath, err)
        } else {
            continue
            //fmt.Printf("Deleted %s\n", filePath)
        }
    }
}