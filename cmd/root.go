package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "go_fractals",
	Short: "Build Fractals with Go!",
	Long: "A Command Line Application to build fractals in Go.",
}

func init() {
	RootCmd.PersistentFlags().IntVarP(&width, "width", "W", 1000, "[OPTIONAL] Set width")
	RootCmd.PersistentFlags().IntVarP(&height, "height", "H", 1000, "[OPTIONAL] Set height")
	RootCmd.PersistentFlags().IntVarP(&fps, "fps", "f", 10, "[OPTIONAL] The framerate of the video.")
	RootCmd.PersistentFlags().StringVarP(&fileName, "fileName", "F", "fractalOutput", "[OPTIONAL] Set file name")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
	  fmt.Fprintln(os.Stderr, err)
	  os.Exit(1)
	}
}