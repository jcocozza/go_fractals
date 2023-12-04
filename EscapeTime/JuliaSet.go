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

type JuliaSet struct {
	Transformation utils.OneParamEquation
	EscapeCondition escapeCondition
	ColorGenerator colorGenerator
	Center complex128
	MaxItr int
	Zoom float64
}

func NewJuliaSet(transform utils.OneParamEquation, escape escapeCondition, colorGen colorGenerator, center complex128, maxItr int, zoom float64) *JuliaSet {
	return &JuliaSet{
		Transformation: transform,
		EscapeCondition: escape,
		ColorGenerator: colorGen,
		Center: center,
		MaxItr: maxItr,
		Zoom: zoom,
	}
}

// calculate the escape time of a given complex number
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

// draw the julia set
func (s *JuliaSet) Draw(path string, width,height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0,0,width,height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			z := complex(float64(x-width/2)/float64(width)*s.Zoom, float64(y-height/2)/float64(height)*s.Zoom)
			escapeTime := s.CalcEscapeTime(z)

			col := s.ColorGenerator(escapeTime)
			img.Set(x,y,col)
		}
	}
	//DrawText(img, "Your Text Here", width-150, 20)
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

// only draw points whose escape time is the max Iteration value
func (s *JuliaSet) DrawFiltered(path string, width,height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0,0,width,height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			z := complex(float64(x-width/2)/float64(width)*s.Zoom, float64(y-height/2)/float64(height)*s.Zoom)
			escapeTime := s.CalcEscapeTime(z)

			if escapeTime == s.MaxItr {
				col := s.ColorGenerator(escapeTime)
				img.Set(x,y,col)
			}
		}
	}

	/*
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
	*/
	return img
}

func newTransform(funcClass utils.TwoParamEquation, varC complex128) utils.OneParamEquation {
	newT := func(z complex128) complex128 {
		return funcClass(z,varC)
	}
	return newT
}

// return a list of Julia Sets
func JuliaEvolution(functionClass utils.TwoParamEquation, cInit, cIncrement complex128, numIncremenets int, center complex128,maxItr int, zoom float64) []*JuliaSet {
	varyingC := cInit

	var juliaSetList []*JuliaSet
	for i := 0; i < numIncremenets; i++ {
		currTransformation := newTransform(functionClass, varyingC)

		escapeCondition := func(z complex128) bool {
			return cmplx.Abs(z) > 2
		}
		tmpJulia := &JuliaSet{currTransformation, escapeCondition, GreyScaleClear, center,maxItr, zoom}
		juliaSetList = append(juliaSetList, tmpJulia)

		varyingC += cIncrement
		utils.ProgressBar(i,numIncremenets)
	}
	return juliaSetList
}

// 2d evolution through parameter space
func EvolveVideo(functionClass func(complex128,complex128) complex128, cInit, cIncrement complex128, numIncrements int, fps int, outputPath string, center complex128, maxItr int, zoom float64, width, height int) {
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

			tmpJulia := JuliaSet{currTransformation, escapeCondition, GreyScale, center ,maxItr, zoom}

			filename := dir + fmt.Sprintf("/image%d.png", i)
			tmpJulia.Draw(filename, width, height)
		}(i, varyingC)

		varyingC += cIncrement
	}

	wg.Wait() // Wait for all goroutines to finish

	inputPattern := dir+"/image%01d.png"
	outputVideo := outputPath

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