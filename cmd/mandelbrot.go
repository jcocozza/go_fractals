package cmd

import (
	"math/cmplx"

	et "github.com/jcocozza/go_fractals/EscapeTime"
	"github.com/jcocozza/go_fractals/utils"

	"github.com/spf13/cobra"
)

var mandelbrotCommand = &cobra.Command{
	Use: "mandelbrot",
	Short: "create a mandelbrot set",
	Long: "Pass in a complex function for the mandelbrot set",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		downloadsPath := utils.GetDownloadDir()

		var mbs et.MandelbrotSet
		if colored {
			mbs = et.MandelbrotSet{
				Transformation: utils.CreateTwoParamEquation(mandelbrotEquation),
				EscapeCondition: func(z complex128) bool {
					return cmplx.Abs(z) > 2
				},
				InitPoint: complex(0,0),
				Center: utils.ParseComplexString(centerPointString),
				ColorGenerator: et.InfernoColor,
				MaxItr: maxItr,
				Zoom: zoom,
			}
		} else {
			mbs = et.MandelbrotSet{
				Transformation: utils.CreateTwoParamEquation(mandelbrotEquation),
				EscapeCondition: func(z complex128) bool {
					return cmplx.Abs(z) > 2
				},
				InitPoint: complex(0,0),
				Center: utils.ParseComplexString(centerPointString),
				ColorGenerator: et.GreyScale,
				MaxItr: maxItr,
				Zoom: zoom,
			}
		}
		mbs.Draw(downloadsPath+"/"+fileName+".png", width, height)
	},
}

func init() {
	rootCmd.AddCommand(mandelbrotCommand)

	mandelbrotCommand.Flags().StringVarP(&mandelbrotEquation, "equation", "e", "", "[REQUIRED] The equation for your mandelbrot set")
	mandelbrotCommand.Flags().BoolVarP(&colored, "color","c", false, "[OPTIONAL] Default Grey Scale")
	mandelbrotCommand.Flags().StringVarP(&fileName, "fileName", "F", "", "[REQUIRED] The file name in the downloads folder")

	mandelbrotCommand.Flags().StringVarP(&centerPointString, "centerPoint","p","0+0i", "[Optional] Set the center point for the fractal")
	mandelbrotCommand.Flags().Float64VarP(&zoom, "zoom","z",4,"[OPTIONAL] Set the zoom; smaller value is more zoomed in")

	mandelbrotCommand.Flags().IntVarP(&maxItr, "maxItr","m",1000,"[OPTIONAL] Set max iterations for time escape")

	mandelbrotCommand.MarkFlagRequired("equation")
	mandelbrotCommand.MarkFlagRequired("fileName")
}
