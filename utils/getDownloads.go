package utils

import (
	"log/slog"
	"os"
	"path/filepath"
)

// get the downloads directory
func GetDownloadDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Error getting user's home directory:", err)
		panic(err)
	}
	// Construct the path to the Downloads folder
	downloadsPath := filepath.Join(homeDir, "Downloads")

	return downloadsPath
}
