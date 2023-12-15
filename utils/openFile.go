package utils

import (
	"log/slog"
	"os/exec"
)

func Open(filePath string) {
	cmd := exec.Command("open", filePath)

	out, err := cmd.CombinedOutput()

	if err != nil {
		slog.Error("Unable to open file: " + filePath, err)
		slog.Error("Error: " + string(out))
		return
	}
}