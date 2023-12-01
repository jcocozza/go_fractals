package cmd

import (
	"fmt"
	"math/cmplx"
	"os"
	"path/filepath"

	et "github.com/jcocozza/go_fractals/EscapeTime"
	"github.com/jcocozza/go_fractals/utils"

	"github.com/spf13/cobra"
)

var juliaEquation string
var colored bool


var juliaCommand = &cobra.Command{
	Use: "julia",
	Short: "create a julia set",
	Long: "Pass in a complex function for the julia set",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user's home directory:", err)
			return
		}
		// Construct the path to the Downloads folder
		downloadsPath := filepath.Join(homeDir, "Downloads")


		var js et.JuliaSet
		if colored {
			js = et.JuliaSet{
				Transformation: utils.ParseEquation(juliaEquation),
				EscapeCondition: func(c complex128) bool {
					return cmplx.Abs(c) > 2
				},
				ColorGenerator: et.GenerateColor,
				MaxItr: 100,
				Zoom: 4,
			}
		} else {
			js = et.JuliaSet{
				Transformation: utils.ParseEquation(juliaEquation),
				EscapeCondition: func(c complex128) bool {
					return cmplx.Abs(c) > 2
				},
				ColorGenerator: et.GreyScale,
				MaxItr: 100,
				Zoom: 4,
			}
		}


		js.Draw(downloadsPath+"/julia.png")

	},
}
func init() {
	rootCmd.AddCommand(juliaCommand)
	juliaCommand.Flags().StringVarP(&juliaEquation, "equation", "e", "", "[REQUIRED] The equation for your julia set")
	juliaCommand.Flags().BoolVarP(&colored, "color","c", false, "[OPTIONAL] Default Grey Scale")
}