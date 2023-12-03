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
	// TODO NOT IMPLEMENTED YET
	rootCmd.PersistentFlags().IntVarP(&width, "width", "w", 1000, "[OPTIONAL] Set width")
	rootCmd.PersistentFlags().IntVarP(&height, "height", "h", 1000, "[OPTIONAL] Set height")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
	  fmt.Fprintln(os.Stderr, err)
	  os.Exit(1)
	}
}