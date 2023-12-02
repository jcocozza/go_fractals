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
var fileName string

var threeDimensional bool
var cInitString string
var cIncrementString string
var numIncrements int
var fps int

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

		js.Draw(downloadsPath+"/"+fileName+".png")
	},
}

var juliaEvolveCommand = &cobra.Command{
	Use: "julia-evolve",
	Short: "evolve a julia set through a parameter",
	Long: "Create a video or 3d stl of the julia set evolving through parameter space",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user's home directory:", err)
			return
		}
		// Construct the path to the Downloads folder
		downloadsPath := filepath.Join(homeDir, "Downloads")

		juliaSetlist := et.JuliaEvolution(
			utils.ParseEquation2(juliaEquation),
			utils.ParseComplexString(cInitString),
			utils.ParseComplexString(cIncrementString),
			numIncrements,
		)
		if threeDimensional {
			stlFileName := fileName+".stl"
			stlFile, err := os.Create(stlFileName)
			if err != nil {
				fmt.Println("Error creating STL file:", err)
				return
			}
			defer stlFile.Close()

			shift := 0.0
			stlFile.WriteString("solid GeneratedModel\n")
			for i,js := range juliaSetlist {
				fmt.Println(js)
				fName := fmt.Sprintf("tmp-%d.png",i)
				et.DrawJuliaSet3D(js.DrawFiltered(fName), stlFile, shift)
				shift += 10
				utils.ProgressBar(i, len(juliaSetlist))
			}
			stlFile.WriteString("endsolid GeneratedModel\n")
		} else {
			et.EvolveVideo(
				utils.ParseEquation2(juliaEquation),
				utils.ParseComplexString(cInitString),
				utils.ParseComplexString(cIncrementString),
				numIncrements,
				fps,
				downloadsPath + "/" + fileName + ".mp4",
			)
		}
	},
}


func init() {
	rootCmd.AddCommand(juliaCommand)
	rootCmd.AddCommand(juliaEvolveCommand)

	// regualar julia set
	juliaCommand.Flags().StringVarP(&juliaEquation, "equation", "e", "", "[REQUIRED] The equation for your julia set")
	juliaCommand.Flags().BoolVarP(&colored, "color","c", false, "[OPTIONAL] Default Grey Scale")
	juliaCommand.Flags().StringVarP(&fileName, "fileName", "F", "", "[REQUIRED] The file name in the downloads folder")
	juliaCommand.MarkFlagRequired("equation")
	juliaCommand.MarkFlagRequired("fileName")

	// julia evolution flags
	juliaEvolveCommand.Flags().StringVarP(&juliaEquation, "equation", "e", "", "[REQUIRED] The parameterized equation for your julia set")
	juliaEvolveCommand.Flags().StringVarP(&cInitString, "initialComplex","p","", "[REQUIRED] Set the intial parameter for a julia evolution")
	juliaEvolveCommand.Flags().StringVarP(&cIncrementString, "complexIncrement","i","", "[REQUIRED] Set the increment for the evolution of the parameter")
	juliaEvolveCommand.Flags().IntVarP(&numIncrements, "numIncrements", "n",10,"[REQUIRED] the number of evolution steps to take")
	juliaEvolveCommand.Flags().StringVarP(&fileName, "fileName", "F", "", "[REQUIRED] The file name in the downloads folder")

	juliaEvolveCommand.MarkFlagsRequiredTogether("equation", "initialComplex", "complexIncrement", "numIncrements")

	juliaEvolveCommand.Flags().BoolVarP(&threeDimensional, "threeDim","d", false, "[OPTIONAL] Create a 3d stl file of the evolution")
	juliaEvolveCommand.Flags().IntVarP(&fps, "fps", "f", 10, "[OPTIONAL] The framerate of the video.")
}