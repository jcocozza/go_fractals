package EscapeTime

import (
	"image"
	"image/png"
	"os"

	"github.com/jcocozza/go_fractals/utils"
)

type MandelbrotSet struct {
	Transformation utils.TwoParamEquation
	InitPoint complex128
	Center complex128
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

func (s *MandelbrotSet) Draw(path string) {
	img := image.NewRGBA(image.Rect(0,0,width,height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			//z := complex(float64(x-width/2)/width*s.Zoom, float64(y-height/2)/height*s.Zoom)

			z := complex(
                float64(x-width/2)/width*s.Zoom+real(s.Center),
                float64(y-height/2)/height*s.Zoom+imag(s.Center),
            )
			escapeTime := s.CalcEscapeTime(z)
			col := s.ColorGenerator(escapeTime) //GenerateColor(escapeTime)
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