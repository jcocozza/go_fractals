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

func NewJuliaSet(transform juliaTransformation, escape escapeCondition, colorGen colorGenerator, maxItr int, zoom float64) *JuliaSet {
	return &JuliaSet{
		Transformation: transform,
		EscapeCondition: escape,
		ColorGenerator: colorGen,
		MaxItr: maxItr,
		Zoom: zoom,
	}
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

func (s *JuliaSet) Draw(path string) *image.RGBA {
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
	return img
}

func (s *JuliaSet) DrawFiltered(path string) *image.RGBA {
	img := image.NewRGBA(image.Rect(0,0,width,height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			z := complex(float64(x-width/2)/width*s.Zoom, float64(y-height/2)/height*s.Zoom)
			escapeTime := s.CalcEscapeTime(z)

			if escapeTime == s.MaxItr {
				col := s.ColorGenerator(escapeTime)
				img.Set(x,y,col)
			}
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
	return img
}

func newTransform(funcClass func(complex128,complex128) complex128, varC complex128) func(complex128) complex128 {
	newT := func(z complex128) complex128 {
		return funcClass(z,varC)
	}
	return newT
}

// return a list of Julia Sets
func JuliaEvolution(functionClass func(complex128,complex128) complex128, cInit, cIncrement complex128, numIncremenets int) []*JuliaSet {
	maxItr := 1000
	zoom := 4.0
	varyingC := cInit

	var juliaSetList []*JuliaSet
	for i := 0; i < numIncremenets; i++ {
		//currTransformation := func(z complex128) complex128 {
		//	return functionClass(z, varyingC)
		//}
		currTransformation := newTransform(functionClass, varyingC)

		escapeCondition := func(z complex128) bool {
			return cmplx.Abs(z) > 2
		}
		tmpJulia := &JuliaSet{currTransformation, escapeCondition, GreyScaleClear, maxItr, zoom}
		juliaSetList = append(juliaSetList, tmpJulia)

		varyingC += cIncrement
		utils.ProgressBar(i,numIncremenets)
	}
	return juliaSetList
}

func EvolveVideo(functionClass func(complex128,complex128) complex128, cInit, cIncrement complex128, numIncrements int, fps int) {
	dir, _ := os.MkdirTemp("","video")
	defer os.RemoveAll(dir)

	varyingC := cInit
	var wg sync.WaitGroup
	for i := 0; i < numIncrements; i++ {
		wg.Add(1) // Increment the WaitGroup counter for each goroutine

		go func(i int, varC complex128) {
			defer wg.Done() // Decrement the WaitGroup counter when the goroutine is done
			defer utils.ProgressBar(i, numIncrements)

			currTransformation := func(z complex128) complex128 {
				return functionClass(z, varC)
			}
			escapeCondition := func(z complex128) bool {
				return cmplx.Abs(z) > 2
			}

			maxItr := 1000
			zoom := 4.0
			tmpJulia := JuliaSet{currTransformation, escapeCondition, GreyScale, maxItr, zoom}

			filename := dir + fmt.Sprintf("/image%d.png", i)
			tmpJulia.Draw(filename)
		}(i, varyingC)

		varyingC += cIncrement
	}

	wg.Wait() // Wait for all goroutines to finish

	inputPattern := dir+"/image%01d.png"
	outputVideo := "TEST.mp4"

    cmd := exec.Command("ffmpeg",
        "-framerate", fmt.Sprint(fps),            // Frame rate
        "-i", inputPattern,           			  // Input image pattern
        "-c:v", "libx264",            			  // Video codec
        "-pix_fmt", "yuv420p",        			  // Pixel format
        outputVideo)

	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error running ffmpeg command:", err)
		fmt.Println("Combined Output:", string(out))
		return
	}
	utils.DeleteFiles("video", "imag*.png") // clean up images afterwards
}


func TestJulia() {
	js := JuliaSet{
		Transformation: func(c complex128) complex128 {
			return c*c - .2 // 1/(c*c - .72i)//
		},
		EscapeCondition: func(c complex128) bool {
			return cmplx.Abs(c) > 2
		},
		ColorGenerator: UnitaryScale,
		MaxItr: 100,
		Zoom: 4,
	}
	js.Draw("juliaCLEAR.png")
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

func TestJuliaEvolve() {

	transformation := func(z,c complex128) complex128 {
		return 1/(z*z + c)
	}
	EvolveVideo(transformation,complex(0,-.63),complex(0, -.001), 100, 10)
}

