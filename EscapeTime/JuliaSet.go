package EscapeTime

import (
	"fmt"
	"image"
	"image/png"
	"math/cmplx"
	"os"
	"os/exec"
	"sync"

	"github.com/jcocozza/go_fractals/utils"
)

// function f(z): C -> C
type juliaTransformation func(complex128) complex128


type JuliaSet struct {
	Transformation juliaTransformation
	EscapeCondition escapeCondition
	ColorGenerator colorGenerator
	MaxItr int
	Zoom float64
}

func (s *JuliaSet) CalcEscapeTime(z complex128) int {
	orbitItrValue := z
	for i := 0; i < s.MaxItr; i++ {
		orbitItrValue = s.Transformation(orbitItrValue)
		if s.EscapeCondition(orbitItrValue) {
			return i
		}
	}
	return s.MaxItr
}

func (s *JuliaSet) Draw(path string) {
	img := image.NewRGBA(image.Rect(0,0,width,height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			z := complex(float64(x-width/2)/width*s.Zoom, float64(y-height/2)/height*s.Zoom)
			escapeTime := s.CalcEscapeTime(z)
			col := s.ColorGenerator(escapeTime)
			img.Set(x,y,col)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

func TestJulia() {
	js := JuliaSet{
		Transformation: func(c complex128) complex128 {
			return 1/(c*c - .72i)//c*c - .2//
		},
		EscapeCondition: func(c complex128) bool {
			return cmplx.Abs(c) > 2
		},
		ColorGenerator: GreyScale,
		MaxItr: 100,
		Zoom: 4,
	}
	js.Draw("juliaGREYTEST.png")
}


func TestJulia2() {
	js := JuliaSet{
		Transformation: utils.ParseEquation("1/(c*c - .72i)"),
		EscapeCondition: func(c complex128) bool {
			return cmplx.Abs(c) > 2
		},
		ColorGenerator: GreyScale,
		MaxItr: 100,
		Zoom: 4,
	}
	js.Draw("funcTest.png")

}


func Evolve(c, cIncrement complex128, numIncrements int, fps int) {
	dir, _ := os.MkdirTemp("","video")
	defer os.RemoveAll(dir)

	varC := c
	var wg sync.WaitGroup
	for i := 0; i < numIncrements; i++ {
		wg.Add(1) // Increment the WaitGroup counter for each goroutine
		//utils.ProgressBar(i, numIncrements)
		//transformation := func(z complex128) complex128 {
		//	return (z*z - varC)
		//}
		go func(i int, varC complex128) {
			defer wg.Done() // Decrement the WaitGroup counter when the goroutine is done
			defer utils.ProgressBar(i, numIncrements)

			transformation := func(c complex128) complex128 {
				return 1/(c*c + varC)
			}
			escapeCondition := func(z complex128) bool {
				return cmplx.Abs(z) > 2
			}

			maxItr := 1000
			zoom := 4.0

			tmpJulia := JuliaSet{transformation, escapeCondition, GenerateColor, maxItr, zoom}

			filename := dir + fmt.Sprintf("/image%d.png", i)
			tmpJulia.Draw(filename)
		}(i, varC)

		varC += cIncrement
	}

	wg.Wait() // Wait for all goroutines to finish

	inputPattern := dir+"/image%01d.png"
	outputVideo := "color.mp4"

    cmd := exec.Command("ffmpeg",
        "-framerate", fmt.Sprint(fps),            // Frame rate
        "-i", inputPattern,           // Input image pattern
        "-c:v", "libx264",            // Video codec
        "-pix_fmt", "yuv420p",        // Pixel format
        outputVideo)

	//err := cmd.Run()
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error running ffmpeg command:", err)
		fmt.Println("Combined Output:", string(out))
		return
	}
	utils.DeleteFiles("video", "imag*.png") // clean up images afterwards
}

func TestJuliaEvolve() {
	Evolve(complex(0,-.63),complex(0, -.001), 100, 10)
}

