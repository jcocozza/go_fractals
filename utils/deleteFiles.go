package utils

import (
	"log/slog"
	"os"
	"path/filepath"
)

func DeleteFiles(dirPath string, filePattern string) {
    // Use filepath.Glob to get a list of matching file paths
    matchingFiles, err := filepath.Glob(filepath.Join(dirPath, filePattern))
    if err != nil {
        slog.Error("Error:", err)
        return
    }

    // Loop through the matching files and delete them
    for _, filePath := range matchingFiles {
        err := os.Remove(filePath)
        if err != nil {
            slog.Error("Error deleting %s: %v\n", filePath, err)
        } else {
            slog.Debug("Deleted " + filePath)
            continue
        }
    }
}