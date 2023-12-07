package cmd

import (
	"fmt"
	"image"
	"math/cmplx"
	"sync"

	et "github.com/jcocozza/go_fractals/EscapeTime"
	"github.com/jcocozza/go_fractals/utils"
	IMGS "github.com/jcocozza/go_fractals/images"

	"github.com/spf13/cobra"
)

var juliaCommand = &cobra.Command{
	Use: "julia",
	Short: "create a julia set",
	Long: "Pass in a complex function for the julia set",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		downloadsPath := utils.GetDownloadDir()


		var js et.JuliaSet
		if colored {
			js = et.JuliaSet{
				Transformation: utils.CreateOneParamEquation(juliaEquation),
				EscapeCondition: func(z complex128) bool {
					return cmplx.Abs(z) > 2
				},
				ColorGenerator: et.InfernoColor,
				Center: utils.ParseComplexString(centerPointString),
				MaxItr: maxItr,
				Zoom: zoom,
			}
		} else {
			js = et.JuliaSet{
				Transformation: utils.CreateOneParamEquation(juliaEquation),
				EscapeCondition: func(z complex128) bool {
					return cmplx.Abs(z) > 2
				},
				ColorGenerator: et.GreyScale,
				Center: utils.ParseComplexString(centerPointString),
				MaxItr: maxItr,
				Zoom: zoom,
			}
		}

		js.Draw(downloadsPath+"/"+fileName+".png", width, height)
	},
}

var juliaEvolveCommand = &cobra.Command{
	Use: "julia-evolve",
	Short: "evolve a julia set through a parameter",
	Long: "Create a video or 3d stl of the julia set evolving through parameter space",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		downloadsPath := utils.GetDownloadDir()

		juliaSetlist := et.JuliaEvolution(
			utils.CreateTwoParamEquation(juliaEquation),
			utils.ParseComplexString(cInitString),
			utils.ParseComplexString(cIncrementString),
			numIncrements,
			utils.ParseComplexString(centerPointString),
			maxItr,
			zoom,
		)
		if threeDimensional {
			// creating a list of fractal images in parallel
			imageChan := make(chan struct {
				Index int // to ensure proper indexing of channel
				Image *image.RGBA
			}, len(juliaSetlist))

			var wg sync.WaitGroup
			for i, js := range juliaSetlist {
				wg.Add(1) // Increment the wait group counter

				go func(i int, js *et.JuliaSet) {
					defer wg.Done() // Decrement the wait group counter when the goroutine completes

					fName := fmt.Sprintf("tmp-%d.png", i)
					img := js.DrawFiltered(fName, width, height) // create the images in parallel

					// Send the *image.RGBA instance to the channel
					imageChan <- struct {
						Index int
						Image *image.RGBA
					}{Index: i, Image: img}
				}(i, js)
			}
			// Close the channel once all goroutines are done
			go func() {
				wg.Wait()
				close(imageChan)
			}()

			// Create a slice to store the images in the correct order
			images := make([]struct {
				Index int
				Image *image.RGBA
			}, len(juliaSetlist))

			// Receive images from the channel and store them in the slice
			for i := range imageChan {
				images[i.Index] = i
			}

			var imgList []image.Image
			for _, imgStruct := range images {
				imgList = append(imgList, imgStruct.Image)
			}

			IMGS.STLControlFlow(writeBinary, solid, imgList, fileName)
		} else {
			et.EvolveVideo(
				utils.CreateTwoParamEquation(juliaEquation),
				utils.ParseComplexString(cInitString),
				utils.ParseComplexString(cIncrementString),
				numIncrements,
				fps,
				downloadsPath + "/" + fileName + ".mp4",
				utils.ParseComplexString(centerPointString),
				maxItr,
				zoom,
				width,
				height,
			)
		}
	},
}


func init() {
	RootCmd.AddCommand(juliaCommand)
	RootCmd.AddCommand(juliaEvolveCommand)

	// regualar julia set
	juliaCommand.Flags().StringVarP(&juliaEquation, "equation", "e", "", "[REQUIRED] The equation for your julia set")
	juliaCommand.Flags().BoolVarP(&colored, "color","c", false, "[OPTIONAL] Default Grey Scale")

	juliaCommand.Flags().Float64VarP(&zoom, "zoom","z",4,"[OPTIONAL] Set the zoom; smaller value is more zoomed in")

	juliaCommand.Flags().IntVarP(&maxItr, "maxItr","m",1000,"[OPTIONAL] Set max iterations for time escape")

	juliaCommand.Flags().StringVarP(&centerPointString, "centerPoint","p","0+0i", "[Optional] Set the center point for the fractal")

	juliaCommand.MarkFlagRequired("equation")

	// julia evolution flags
	juliaEvolveCommand.Flags().StringVarP(&juliaEquation, "equation", "e", "", "[REQUIRED] The parameterized equation for your julia set")
	juliaEvolveCommand.Flags().StringVarP(&cInitString, "initialComplex","P","", "[REQUIRED] Set the intial parameter for a julia evolution")
	juliaEvolveCommand.Flags().StringVarP(&cIncrementString, "complexIncrement","i","", "[REQUIRED] Set the increment for the evolution of the parameter")
	juliaEvolveCommand.Flags().IntVarP(&numIncrements, "numIncrements", "n",10,"[REQUIRED] the number of evolution steps to take")

	juliaEvolveCommand.MarkFlagsRequiredTogether("equation", "initialComplex", "complexIncrement", "numIncrements")

	juliaEvolveCommand.Flags().BoolVarP(&threeDimensional, "threeDim","d", false, "[OPTIONAL] Create a 3d stl file of the evolution")
	juliaEvolveCommand.Flags().BoolVarP(&writeBinary, "binary","b", false, "[OPTIONAL] save the stl as a binary")
	juliaEvolveCommand.Flags().BoolVarP(&solid, "solid","s", false, "[OPTIONAL] write the stl as a completely solid object(much larger file size)")
	juliaEvolveCommand.Flags().IntVarP(&fps, "fps", "f", 10, "[OPTIONAL] The framerate of the video.")
	juliaEvolveCommand.Flags().Float64VarP(&zoom, "zoom","z",4,"[OPTIONAL] Set the zoom; smaller value is more zoomed in")

	juliaEvolveCommand.Flags().IntVarP(&maxItr, "maxItr","m",1000,"[OPTIONAL] Set max iterations for time escape")
	juliaEvolveCommand.Flags().StringVarP(&centerPointString, "centerPoint","p","0+0i", "[OPTIONAL] Set the center point for the fractal")
}