package EscapeTime

import (
	"image"
	"image/png"
	"math/cmplx"
	"os"
)

// function f(z,c): C,C -> C
type mandelbrotTransformation func(complex128, complex128) complex128

type MandelbrotSet struct {
	Transformation mandelbrotTransformation
	InitPoint complex128
	EscapeCondition escapeCondition
	ColorGenerator colorGenerator
	MaxItr int
	Zoom float64
}

func (s *MandelbrotSet) CalcEscapeTime(c complex128) int {
	orbitItrValue := s.InitPoint
	for i := 0; i < s.MaxItr; i++ {
		orbitItrValue = s.Transformation(orbitItrValue, c)
		if s.EscapeCondition(orbitItrValue) {
			return i
		}
	}
	return s.MaxItr
}

func (s *MandelbrotSet) Draw() {
	img := image.NewRGBA(image.Rect(0,0,width,height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			z := complex(float64(x-width/2)/width*s.Zoom, float64(y-height/2)/height*s.Zoom)
			//z := complex(float64(x-width/4)/width*s.Zoom, float64(y-height/4)/height*s.Zoom)

			escapeTime := s.CalcEscapeTime(z)
			col := s.ColorGenerator(escapeTime) //GenerateColor(escapeTime)
			img.Set(x,y,col)
		}
	}

	file, err := os.Create("mandelbrot.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

func TestMandelbrot() {
	mb := MandelbrotSet{
		Transformation: func(z, c complex128) complex128 {
			return z*z + c //return 1/(z*z + c)
		},
		InitPoint: complex(0,0),
		EscapeCondition: func(c complex128) bool {
			return cmplx.Abs(c) > 2
		},
		ColorGenerator: GreyScale,
		MaxItr: 1000,
		Zoom: 4,
	}
	mb.Draw()
}