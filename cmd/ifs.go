package cmd

import (
	"fmt"
	"strings"

	IFS "github.com/jcocozza/go_fractals/IteratedFunctionSystems"
	"github.com/jcocozza/go_fractals/utils"
	viz "github.com/jcocozza/go_fractals/visualizer"
	"github.com/spf13/cobra"
	"gonum.org/v1/gonum/mat"
	BAR "github.com/schollz/progressbar/v3"

)
var Path string
var numIterations int
var createVideo bool
var algorithmProbabilistic bool
var algorithmDeterministic bool
var frameRate int

/* // TODO
var out string - let user specify where the image/video writes to
var width, height int - let user specify width, height of image/video
var framerate int - let user specify framerate of the video
*/

func init() {
	rootCmd.AddCommand(ifsCmd)

	ifsCmd.Flags().BoolVar(&algorithmProbabilistic, "algo-p", false, "[OPTIONAL] Use the probabilistic algorithm")
	ifsCmd.Flags().BoolVar(&algorithmDeterministic, "algo-d", false, "[OPTIONAL] Use the deterministic algorithm")
	ifsCmd.MarkFlagsMutuallyExclusive("algo-p", "algo-d")
	ifsCmd.MarkFlagRequired("algorithm-probabilistic")
	ifsCmd.MarkFlagRequired("algorithm-deterministic")

	ifsCmd.Flags().StringVarP(&Path, "path", "p", "", "[REQUIRED] The path to your iterated function system file")
	ifsCmd.MarkFlagRequired("path")

	ifsCmd.Flags().IntVarP(&numIterations, "numItr", "n", 1, "[OPTIONAL] The number of iterations you want to use. Deterministic algorithm: relatively low because of exponential growth.")

	ifsCmd.Flags().BoolVarP(&createVideo, "video","v", false, "[OPTIONAL] Whether to create a video or not")

	ifsCmd.Flags().IntVarP(&frameRate, "fps", "f", 10, "[OPTIONAL] The framerate of the video.")
}

var ifsCmd = &cobra.Command{
	Use: "ifs",
	Short: "Run an iterated function system",
	Long: "Pass in a file that contains an iterated function system",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		bar := BAR.Default(int64(numIterations + 1), "Generating Fractals...")
		//fmt.Println("Reading from file:", Path)

		rowMatches := utils.Reader(Path, `\[(.*?)\]`)

		var transformationList []IFS.Transformation
		var dimension int
		for _, rowMatch := range rowMatches {
			dimSimilarity := utils.ParseCommaDelimitedStr(strings.Trim(rowMatch[0], "[]"))
			similarity := utils.ParseCommaDelimitedStr(strings.Trim(rowMatch[1], "[]"))

			dimShift := utils.ParseCommaDelimitedStr(strings.Trim(rowMatch[2], "[]"))
			shift := utils.ParseCommaDelimitedStr(strings.Trim(rowMatch[3], "[]"))

			dimension = int(dimShift[0])

			similarityMatrix := mat.NewDense(int(dimSimilarity[0]),int(dimSimilarity[1]), similarity)
			shiftMatrix := mat.NewDense(int(dimShift[0]),int(dimShift[1]), shift)
			transform, err := IFS.NewTransformation(*similarityMatrix, *shiftMatrix)

			if err != nil {
				panic(fmt.Sprintf("Unable to create transform: %s", err))
			}

			transformationList = append(transformationList, *transform)
		}

		newIfs := IFS.NewIteratedFunctionSystem(transformationList, numIterations, dimension)
		const width, height = 1000,1000

		if dimension == 3 {
			if algorithmDeterministic {
				pointsList := newIfs.RunDeterministic()
				viz.Draw3D(pointsList)
				return
			} else if algorithmProbabilistic {
				pointsList := newIfs.RunProbabilistic()
				viz.Draw3D(pointsList)
				return
			} else {
				panic("No other algorithm to use!!")
			}

		} else if dimension == 2 {

			if algorithmProbabilistic {
				if createVideo {
					progressCh := make(chan int)
					go viz.VideoWrapper(width, height, "my_fractal_video.mp4", *newIfs, newIfs.RunProbabilisticStepWise, frameRate, progressCh)

					// Monitor the progress channel and update the progress bar
					for progress := range progressCh {
						bar.Set(progress)
					}
					// Finish the progress bar when the Goroutine is done
					bar.Finish()
					return
				}
				pointsList := newIfs.RunProbabilistic()
				fractal := viz.NewFractalImage(width, height, "my_fractal.png", pointsList)
				fractal.WriteImage()
			} else if algorithmDeterministic {
				if createVideo {
					progressCh := make(chan int)
					go viz.VideoWrapper(width, height, "my_fractal_video.mp4", *newIfs, newIfs.RunDeterministicStepWise, frameRate, progressCh)
					// Monitor the progress channel and update the progress bar
					for progress := range progressCh {
						bar.Set(progress)
					}
					// Finish the progress bar when the Goroutine is done
					bar.Finish()
					return
				}
				pointsList := newIfs.RunDeterministic()
				fractal := viz.NewFractalImage(width, height, "my_fractal.png", pointsList)
				fractal.WriteImage()
			} else {
				panic("No other algorithm to use!!")
			}
		}
	},
}