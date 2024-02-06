package gui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/cmplx"
	"os"

	"github.com/jcocozza/go_fractals/EscapeTime"
	"github.com/jcocozza/go_fractals/images"
	"github.com/jcocozza/go_fractals/utils"
)

func CreateJuliaSets(cList []complex128) []*EscapeTime.JuliaSet {
	jsList := []*EscapeTime.JuliaSet{}
	for _, pt := range cList {
		currentPt := pt // Create a local variable
		js := EscapeTime.NewJuliaSet(
			func(z complex128) complex128 { return z*z + currentPt },
			func(z complex128) bool { return cmplx.Abs(z) > 2 },
			EscapeTime.GreyScale,
			complex(0,0),
			1000,
			4,
		)
		jsList = append(jsList, js)
	}
	return jsList
}

func CreateImagePairs(jsList []*EscapeTime.JuliaSet, ptList []*image.Point, mbs EscapeTime.MandelbrotSet) []*image.RGBA {
	imgPairlist := []*image.RGBA{}
	mbimg := mbs.DrawImg(800,800)
	for i, js := range jsList {
		currPt := ptList[i]
		jsImg := js.Draw("", 800, 800)
		mbimg.Set(currPt.X, currPt.Y, color.RGBA{0,255,0,255})

		// Create a new image with the combined width and the maximum height of the input images
		resultImg := image.NewRGBA(image.Rect(0, 0, jsImg.Bounds().Dx()+mbimg.Bounds().Dx(), max(jsImg.Bounds().Dy(), mbimg.Bounds().Dy())))

		// Draw the first image onto the result image at the left side
		draw.Draw(resultImg, jsImg.Bounds(), jsImg, image.Point{}, draw.Over)
		// Draw the second image onto the result image next to the first image
		draw.Draw(resultImg, mbimg.Bounds().Add(image.Pt(jsImg.Bounds().Dx(), 0)), mbimg, image.Point{}, draw.Over)

		imgPairlist = append(imgPairlist, resultImg)
		utils.ProgressBar(i, len(jsList))
	}
	return imgPairlist
}

func CreateVideo(imgList []*image.RGBA) {
	dir, _ := os.MkdirTemp("","video")
	defer os.RemoveAll(dir)

	for i,img := range imgList {
		filename := dir + fmt.Sprintf("/image%d.png", i)
		images.SavePNG(img, filename)
	}

	inputPattern := dir+"/image%01d.png"
	images.CreateVideo(inputPattern, "/Users/josephcocozza/Downloads/OUTPUT.mp4",10)
}

