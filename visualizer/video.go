package visualizer

import (
	"fmt"
	IFS "github.com/jcocozza/go_fractals/IteratedFunctionSystems"
	"os/exec"
)

type FractalVideo struct {
	Width int
	Height int
	Path string
	IFSys IFS.IteratedFunctionSystem
}

func NewFractalVideo(width int, height int, path string, ifs IFS.IteratedFunctionSystem) *FractalVideo {
	return &FractalVideo{
		Width: width,
		Height: height,
		Path: path,
		IFSys: ifs,
	}
}

//WriteVideoImages will run the deterministic algorithm one step at a time and save an image at each step
func (fv *FractalVideo) WriteVideoImages(numIterations int) {
	pointsList := fv.IFSys.InitialPoints

	// write the initial conditions
	fractal := NewFractalImage(fv.Width, fv.Height, "video/image0.png", pointsList)
	fractal.WriteImage()

	for i := 0; i < numIterations; i ++ {
		pointsList = fv.IFSys.RunDeterministicStepWise(pointsList)
		fractal := NewFractalImage(fv.Width, fv.Height, fmt.Sprintf("video/image%d.png", i), pointsList)
		fractal.WriteImage()
	}
}

//CreateVideo combine the images into a video
func (fv *FractalVideo) CreateVideo() {
	inputPattern := "video/image%01d.png"
	outputVideo := fv.Path

    cmd := exec.Command("ffmpeg",
        "-framerate", "1",            // Frame rate
        "-i", inputPattern,           // Input image pattern
        "-c:v", "libx264",            // Video codec
        "-pix_fmt", "yuv420p",        // Pixel format
        outputVideo)

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Fractal Video Created")
	}
}