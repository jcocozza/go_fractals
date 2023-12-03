package cmd

import (
	"os"
	"path/filepath"

	IFS "github.com/jcocozza/go_fractals/IteratedFunctionSystems"
	"github.com/jcocozza/go_fractals/stack"
	viz "github.com/jcocozza/go_fractals/visualizer"
	BAR "github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

/* // TODO
var out string - let user specify where the image/video writes to
var width, height int - let user specify width, height of image/video
*/

func init() {
	rootCmd.AddCommand(ifsCmd)

	ifsCmd.Flags().BoolVar(&algorithmProbabilistic, "algo-p", false, "[OPTIONAL] Use the probabilistic algorithm")
	ifsCmd.Flags().BoolVar(&algorithmDeterministic, "algo-d", true, "[OPTIONAL] Use the deterministic algorithm")
	ifsCmd.MarkFlagsMutuallyExclusive("algo-p", "algo-d")
	//ifsCmd.MarkFlagRequired("algorithm-probabilistic")
	//ifsCmd.MarkFlagRequired("algorithm-deterministic")

	ifsCmd.Flags().StringVarP(&Path, "path", "p", "", "[REQUIRED] The path to your iterated function system file")
	//ifsCmd.MarkFlagRequired("path")

	ifsCmd.Flags().IntVarP(&numIterations, "numItr", "n", 1, "[OPTIONAL] The number of iterations you want to use.")

	ifsCmd.Flags().BoolVarP(&createVideo, "video","v", false, "[OPTIONAL] Whether to create a video or not")

	ifsCmd.Flags().IntVarP(&fps, "fps", "f", 10, "[OPTIONAL] The framerate of the video.")

	ifsCmd.Flags().Float64SliceVar(&probabilitiesList, "probabilities", []float64{}, "[OPTIONAL - comma separated] Specify probabilities of transformations. Must add to 1. If none will calculated based on matrices. Note that a determinant of zero can cause unexpected things.")
	ifsCmd.MarkFlagsMutuallyExclusive("algo-d", "probabilities")
	ifsCmd.MarkFlagsMutuallyExclusive("probabilities", "video")

	ifsCmd.Flags().IntVarP(&numPoints, "numPoints", "z", 1, "[OPTIONAL] The number of initial points.")

	ifsCmd.Flags().BoolVarP(&fractalStack, "stack","s", false, "[OPTIONAL] Generate the corresponding fractal stack - writes to ~/Downloads/out.stl file")
	ifsCmd.Flags().IntVarP(&numStacks, "numStacks", "k", 1, "[OPTIONAL] The number of stacks to generate")
	ifsCmd.Flags().Float32VarP(&thickness, "thickness","T", 15, "[OPTIONAL] Specify the thickness the stack layer")
	ifsCmd.MarkFlagsRequiredTogether("stack", "numStacks", "thickness")

	ifsCmd.Flags().BoolVarP(&random, "random","r", false, "[OPTIONAL] Create a random 2D Iterated Function system using the probabilistic algorithm")
	//ifsCmd.MarkFlagsMutuallyExclusive("random", "video", "probabilities", "fps", "algo-p", "algo-d", "path")
	ifsCmd.Flags().IntVarP(&numTransforms, "numTransforms", "t", 2, "[OPTIONAL] The number of transforms to randomly generate.")
	//rootCmd.MarkFlagsRequiredTogether("random", "numTransforms")

}

var ifsCmd = &cobra.Command{
	Use: "ifs",
	Short: "Run an iterated function system",
	Long: "Pass in a file that contains an iterated function system",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		const width, height = 1000,1000
		// Get the user's home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			//fmt.Println("Error getting user's home directory:", err)
			return
		}

		// Construct the path to the Downloads folder
		downloadsPath := filepath.Join(homeDir, "Downloads")


		if random {
			randIFS := IFS.GenerateRandomIFS(2, numTransforms)
			pointsList := randIFS.RunProbabilistic(randIFS.CalculateProbabilities())
			//fmt.Println("POints list", pointsList[1])
			fractal := viz.NewFractalImage(width, height, downloadsPath+"/random_fractal.png", pointsList)
			fractal.WriteImage()
			return
		}

		transformationList, dimension := IFS.ParseIFS(Path)

		if len(probabilitiesList) != 0 && len(probabilitiesList) != len(transformationList) {
			panic("You must pass as many probabilities as there are transforms")
		}

		probSum := 0.0
		for _, prob := range probabilitiesList {
			probSum += prob
		}

		if len(probabilitiesList) != 0 && probSum != 1 {
			panic("Passed probabilities must sum to 1")
		}


		newIfs := IFS.NewIteratedFunctionSystem(transformationList, numIterations, numPoints, dimension)

		if fractalStack {
			if numStacks == 1 { // in the case of a stack with length 1, we really just want to thicken a fractal after it's algorithm has been run
				// choice of which algorithm generates the points in the case of the fractal thickening
				if algorithmProbabilistic {
					if len(probabilitiesList) == 0 {
						pointsList := newIfs.RunProbabilistic(newIfs.CalculateProbabilities())
						fractal := viz.NewFractalImage(width, height, "stack.png", pointsList)
						fractal.WriteImage()
						stack.ThickenedFractal(fractal, thickness, downloadsPath+"/out.stl")
					} else {
						pointsList := newIfs.RunProbabilistic(probabilitiesList)
						fractal := viz.NewFractalImage(width, height, "stack.png", pointsList)
						fractal.WriteImage()
						stack.ThickenedFractal(fractal, thickness, downloadsPath+"/out.stl")
					}

				} else if algorithmDeterministic {
					pointsList := newIfs.RunDeterministic()
					fractal := viz.NewFractalImage(width, height, downloadsPath+"/deterministic_fractal.png", pointsList)
					fractal.WriteImage()
					stack.ThickenedFractal(fractal, thickness, downloadsPath+"/out.stl")
				}
			} else {
				stack.CreateFractalStack(Path, numStacks, thickness, downloadsPath+"/out.stl")
			}
			return
		}

		if dimension == 3 {
			if algorithmDeterministic {
				pointsList := newIfs.RunDeterministic()
				viz.Draw3D(pointsList)
				return
			} else if algorithmProbabilistic {
				if len(probabilitiesList) == 0 {
					pointsList := newIfs.RunProbabilistic(newIfs.CalculateProbabilities())
					viz.Draw3D(pointsList)
					return
				} else {
					pointsList := newIfs.RunProbabilistic(probabilitiesList)
					viz.Draw3D(pointsList)
					return
				}

			} else {
				panic("No other algorithm to use!!")
			}

		} else if dimension == 2 {

			if algorithmProbabilistic {
				if createVideo {
					bar := BAR.Default(int64(numIterations + 1), "Generating Fractals...")
					progressCh := make(chan int)
					go viz.VideoWrapper(width, height, downloadsPath+"/fractal_video.mp4", *newIfs, newIfs.RunProbabilisticStepWise, fps, progressCh)

					// Monitor the progress channel and update the progress bar
					for progress := range progressCh {
						bar.Set(progress)
					}
					// Finish the progress bar when the Goroutine is done
					bar.Finish()
					return
				}

				if len(probabilitiesList) == 0 {
					pointsList := newIfs.RunProbabilistic(newIfs.CalculateProbabilities())
					fractal := viz.NewFractalImage(width, height, downloadsPath+"/probabilistic_fractal.png", pointsList)
					fractal.WriteImage()
				} else {
					pointsList := newIfs.RunProbabilistic(probabilitiesList)
					fractal := viz.NewFractalImage(width, height, downloadsPath+"/probabilistic_fractal.png", pointsList)
					fractal.WriteImage()
				}

			} else if algorithmDeterministic {
				if createVideo {
					bar := BAR.Default(int64(numIterations + 1), "Generating Fractals...")
					progressCh := make(chan int)
					go viz.VideoWrapper(width, height, downloadsPath+"/fractal_video.mp4", *newIfs, newIfs.RunDeterministicStepWise, fps, progressCh)
					// Monitor the progress channel and update the progress bar
					for progress := range progressCh {
						bar.Set(progress)
					}
					// Finish the progress bar when the Goroutine is done
					bar.Finish()
					return
				}
				pointsList := newIfs.RunDeterministic()
				fractal := viz.NewFractalImage(width, height, downloadsPath+"/deterministic_fractal.png", pointsList)
				fractal.WriteImage()
			} else {
				panic("No other algorithm to use!!")
			}
		}
	},
}