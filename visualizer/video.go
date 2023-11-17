package visualizer

import (
	"fmt"
	"os/exec"
	"os"

	IFS "github.com/jcocozza/go_fractals/IteratedFunctionSystems"
	"github.com/jcocozza/go_fractals/utils"
)

type FractalVideo struct {
	Width int
	Height int
	Path string
	IFSys IFS.IteratedFunctionSystem
	FrameRate int
}

// An Algorithm is a function that takes in a list of points and returns a list of points
type Algorithm func([][]float64) [][]float64

func VideoWrapper(width int, height int, fileName string, ifs IFS.IteratedFunctionSystem, stepWiseAlgo Algorithm, frameRate int, progressCh chan int) {
	fv := newFractalVideo(width, height, fileName, ifs, frameRate)

	dir, _ := os.MkdirTemp("","video")
	defer os.RemoveAll(dir)

	pointAccumulator := fv.IFSys.InitialPoints
	newPoints := fv.IFSys.InitialPoints
	// write the initial conditions
	fractal := NewFractalImage(fv.Width, fv.Height, dir+"/image0.png", pointAccumulator)
	fractal.WriteImage()

	for i := 0; i < fv.IFSys.NumIterations; i ++ {
		newPoints = stepWiseAlgo(newPoints)
		pointAccumulator = append(pointAccumulator, newPoints...)
		//fmt.Println("POINT SET:", pointAccumulator)
		fractal := NewFractalImage(fv.Width, fv.Height, fmt.Sprintf(dir+"/image%d.png", i), pointAccumulator)
		fractal.WriteImage()
		progressCh <- (i + 1)
	}
	fv.createVideo(dir)
	progressCh <- fv.IFSys.NumIterations + 1
	close(progressCh)
}

func newFractalVideo(width int, height int, path string, ifs IFS.IteratedFunctionSystem, frameRate int) *FractalVideo {
	return &FractalVideo{
		Width: width,
		Height: height,
		Path: path,
		IFSys: ifs,
		FrameRate: frameRate,
	}
}

//CreateVideo combine the images into a video
func (fv *FractalVideo) createVideo(tmpDir string) {
	inputPattern := tmpDir+"/image%01d.png"
	outputVideo := fv.Path

    cmd := exec.Command("ffmpeg",
        "-framerate", fmt.Sprint(fv.FrameRate),            // Frame rate
        "-i", inputPattern,           // Input image pattern
        "-c:v", "libx264",            // Video codec
        "-pix_fmt", "yuv420p",        // Pixel format
        outputVideo)

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
	//fmt.Println("Fractal Video Created")
	utils.DeleteFiles("video", "imag*.png") // clean up images afterwards
}