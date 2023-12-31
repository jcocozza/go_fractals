package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "go_fractals",
	Short: "Build Fractals with Go!",
	Long: "A Command Line Application to build fractals in Go.",
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&width, "width", "W", 1000, "[OPTIONAL] Set width")
	rootCmd.PersistentFlags().IntVarP(&height, "height", "H", 1000, "[OPTIONAL] Set height")
	rootCmd.PersistentFlags().IntVarP(&fps, "fps", "f", 10, "[OPTIONAL] The framerate of the video.")
	rootCmd.PersistentFlags().StringVarP(&fileName, "fileName", "F", "fractalOutput", "[OPTIONAL] Set file name")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
	  fmt.Fprintln(os.Stderr, err)
	  os.Exit(1)
	}
}