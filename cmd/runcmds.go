package cmd

import (
	"log/slog"

	"github.com/jcocozza/go_fractals/internal/utils"
)

// Run a list of command args
func RunCmdArgs(cmdArgs []string) (string, error) {
	slog.Info("Running command: " + utils.StrListToString(cmdArgs))

	//var outputBuffer bytes.Buffer
	//RootCmd.SetOutput(&outputBuffer)

	RootCmd.SetArgs(cmdArgs)
	if err := RootCmd.Execute(); err != nil {
		slog.Error("error", err)
		return "", err
	}
	//RootCmd.SetOut(os.Stdout)

	return "outputBuffer.String()", nil
}