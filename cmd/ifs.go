package cmd

import (
	"fmt"
	"image"

	IFS "github.com/jcocozza/go_fractals/IteratedFunctionSystems"
	IMGS "github.com/jcocozza/go_fractals/images"
	"github.com/jcocozza/go_fractals/utils"
	"github.com/spf13/cobra"
)

// create a new ifs from user input
func processIFS(path string, probList []float64, numItr int, numPts int) (*IFS.IteratedFunctionSystem, int) {
	// create the transformations
	transformationList, dimension := IFS.ParseIFS(path)

	//check probabilities
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

	// create a new iterated function system
	newIfs := IFS.NewIteratedFunctionSystem(transformationList, numItr, numPts, dimension)
	return newIfs, dimension
}

var ifsCMD = &cobra.Command{
	Use: "ifs",
	Short: "Access Iterated function system commands",
}

var ifsImg = &cobra.Command{
	Use: "image",
	Short: "Run an iterated function system",
	Long: "Pass in a file that contains an iterated function system",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		downloadsPath := utils.GetDownloadDir()
		filePath := downloadsPath + "/" + fileName + ".png"
		if random {
			randIFS := IFS.GenerateRandomIFS(2, numTransforms)
			pointsList := randIFS.RunProbabilistic(randIFS.CalculateProbabilities())
			fractal := IFS.NewFractalImage(width, height, downloadsPath+"/random_fractal.png", pointsList)
			fractal.WriteImage(filePath)
			return
		}

		// create a new iterated function system
		newIfs, dimension := processIFS(ifsPath, probabilitiesList, numIterations, numPoints)

		if dimension != 2 {
			panic("dimensionality of > 2 not currently supported!")
		}

		if algorithmProbabilistic {
			if len(probabilitiesList) == 0 {
				pointsList := newIfs.RunProbabilistic(newIfs.CalculateProbabilities())
				fractal := IFS.NewFractalImage(width, height, downloadsPath+"/probabilistic_fractal.png", pointsList)
				fractal.WriteImage(filePath)
			} else {
				pointsList := newIfs.RunProbabilistic(probabilitiesList)
				fractal := IFS.NewFractalImage(width, height, downloadsPath+"/probabilistic_fractal.png", pointsList)
				fractal.WriteImage(filePath)
			}
			} else if algorithmDeterministic {
				pointsList := newIfs.RunDeterministic()
				fractal := IFS.NewFractalImage(width, height, downloadsPath+"/deterministic_fractal.png", pointsList)
				fractal.WriteImage(filePath)
			} else {
				panic("No other algorithm to use!!")
			}
	},
}

var ifsEvolveCMD = &cobra.Command{
	Use: "evolve",
	Short: "iterate through an ifs",
	Long:  "Create a video or 3d stl of the ifs",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		downloadsPath := utils.GetDownloadDir()
		// create a new iterated function system
		newIfs, _ := processIFS(ifsPath, probabilitiesList, numIterations, numPoints)
		var imgList []image.Image

		if threeDimensional {
			if numStacks == 1 {
				if algorithmProbabilistic {
					if len(probabilitiesList) == 0 {
						pointsList := newIfs.RunProbabilistic(newIfs.CalculateProbabilities())
						fractal := IFS.NewFractalImage(width, height, "stack.png", pointsList)
						imgList = append(imgList, fractal.Img)
					} else {
						pointsList := newIfs.RunProbabilistic(probabilitiesList)
						fractal := IFS.NewFractalImage(width, height, "stack.png", pointsList)
						imgList = append(imgList, fractal.Img)
					}
				} else if algorithmDeterministic {
					pointsList := newIfs.RunDeterministic()
					fractal := IFS.NewFractalImage(width, height, "stack.png", pointsList)
					imgList = append(imgList, fractal.Img)
				}
			} else {
				base := newIfs.InitialPoints
				baseImg := IFS.NewFractalImage(width,height,"base.png", base)

				var tempPoints = newIfs.InitialPoints
				var tempFractalImg *IFS.FractalImage
				var fractalImageList []*IFS.FractalImage

				fractalImageList = append(fractalImageList, baseImg)
				for i := 0; i < numStacks; i++ {
					tempPoints = newIfs.RunDeterministicStepWise(tempPoints)
					tempFractalImg = IFS.NewFractalImage(width,height,fmt.Sprintf("stack-%d",i), tempPoints)
					fractalImageList = append(fractalImageList, tempFractalImg)
				}

				for _, fractal := range fractalImageList {
					imgList = append(imgList, fractal.Img)
				}
			}
			IMGS.STLControlFlow(writeBinary, solid, imgList, fileName)
		} else { // create a video
			pointAccumulator := newIfs.InitialPoints
			newPoints := newIfs.InitialPoints
			// write the initial conditions
			fractal := IFS.NewFractalImage(width, height, "base.png", pointAccumulator)
			imgList = append(imgList, fractal.Img)

			for i := 0; i < newIfs.NumIterations; i ++ {

				if algorithmDeterministic {
					newPoints = newIfs.RunDeterministicStepWise(newPoints)
				} else if algorithmProbabilistic {
					newPoints = newIfs.RunProbabilisticStepWise(newPoints)
				}
				pointAccumulator = append(pointAccumulator, newPoints...)
				fractal := IFS.NewFractalImage(width, height, fmt.Sprintf("/image%d.png", i), pointAccumulator)
				imgList = append(imgList, fractal.Img)
			}

			outPath := downloadsPath + "/" + fileName + ".mp4"
			IMGS.VideoFromImages(imgList, outPath, fps)
		}
	},
}


func init() {
	RootCmd.AddCommand(ifsCMD)

	ifsCMD.AddCommand(ifsImg)
	ifsCMD.AddCommand(ifsEvolveCMD)

	ifsCMD.PersistentFlags().BoolVar(&algorithmProbabilistic, "algo-p", false, "[OPTIONAL] Use the probabilistic algorithm")
	ifsCMD.PersistentFlags().BoolVar(&algorithmDeterministic, "algo-d", true, "[OPTIONAL] Use the deterministic algorithm")
	ifsCMD.MarkFlagsMutuallyExclusive("algo-p", "algo-d")

	ifsCMD.PersistentFlags().StringVarP(&ifsPath, "path", "p", "", "[REQUIRED] The path to your iterated function system file")
	ifsCMD.MarkFlagRequired("path")

	ifsCMD.PersistentFlags().IntVarP(&numIterations, "numItr", "n", 1, "[OPTIONAL] The number of iterations you want to use.")

	ifsCMD.PersistentFlags().Float64SliceVar(&probabilitiesList, "probabilities", []float64{}, "[OPTIONAL - comma separated] Specify probabilities of transformations. Must add to 1. If none will calculated based on matrices. Note that a determinant of zero can cause unexpected things.")
	ifsCMD.MarkFlagsMutuallyExclusive("algo-d", "probabilities")

	ifsCMD.PersistentFlags().IntVarP(&numPoints, "numPoints", "z", 1, "[OPTIONAL] The number of initial points.")

	ifsImg.Flags().BoolVarP(&random, "random","r", false, "[OPTIONAL] Create a random 2D Iterated Function system using the probabilistic algorithm")
	ifsImg.Flags().IntVarP(&numTransforms, "numTransforms", "t", 2, "[OPTIONAL] The number of transforms to randomly generate.")

	ifsEvolveCMD.Flags().BoolVarP(&threeDimensional, "threeDim","d", false, "[OPTIONAL] Create a 3d stl file of the evolution")
	ifsEvolveCMD.PersistentFlags().IntVarP(&numStacks, "numStacks", "k", 1, "[OPTIONAL] The number of stacks to generate")
	ifsEvolveCMD.PersistentFlags().Float32VarP(&thickness, "thickness","T", 15, "[OPTIONAL] Specify the thickness the stack layer")
	ifsEvolveCMD.MarkFlagsRequiredTogether("threeDim","numStacks", "thickness")
}