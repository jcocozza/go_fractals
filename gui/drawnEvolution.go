package gui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/cmplx"
	"os"
	"sync"

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

func setPixelAndReturnCopy(img image.Image, x, y int, col color.RGBA) image.Image {
    // Create a new RGBA image based on the original one
    bounds := img.Bounds()
    newImg := image.NewRGBA(bounds)

    // Draw the original image onto the new one
    draw.Draw(newImg, bounds, img, bounds.Min, draw.Src)

    // Set the specified pixel color in the new image
    newImg.Set(x, y, col)

    return newImg
}

func CreateMandelbrotPath(ptList []*image.Point, mbs EscapeTime.MandelbrotSet) []image.Image {
	mbimg := mbs.DrawImg(800,800)
	var mbTemp image.Image = mbimg
	imgList := []image.Image{}
	for i := range ptList {
		currPt := ptList[i]
		mbTemp = setPixelAndReturnCopy(mbTemp, currPt.X, currPt.Y, color.RGBA{0,255,0,255})
		imgList = append(imgList, mbTemp)
	}
	return imgList
}

func CreateImagePairsParallel(jsList []*EscapeTime.JuliaSet, ptList []*image.Point, mbs EscapeTime.MandelbrotSet)  {
	dir, _ := os.MkdirTemp("","video")
	defer os.RemoveAll(dir)

	fmt.Println(dir)

	//mbimg := mbs.DrawImg(800,800)

	mbPathLst := CreateMandelbrotPath(ptList, mbs)

	var wg sync.WaitGroup
	for i := range jsList {
		wg.Add(1)

		go func (j int)  {
			defer wg.Done()
			utils.ProgressBar(j, len(jsList))
			//currPt := ptList[j]
			jsImg := jsList[j].Draw("", 800, 800)
			newImg := mbPathLst[j]

			// Create a new image with the combined width and the maximum height of the input images
			resultImg := image.NewRGBA(image.Rect(0, 0, jsImg.Bounds().Dx()+newImg.Bounds().Dx(), max(jsImg.Bounds().Dy(), newImg.Bounds().Dy())))

			// Draw the first image onto the result image at the left side
			draw.Draw(resultImg, jsImg.Bounds(), jsImg, image.Point{}, draw.Over)
			// Draw the second image onto the result image next to the first image
			draw.Draw(resultImg, newImg.Bounds().Add(image.Pt(jsImg.Bounds().Dx(), 0)), newImg, image.Point{}, draw.Over)

			filename := dir + fmt.Sprintf("/image%d.png", j)
			images.SavePNG(resultImg, filename)
		}(i)
	}
	wg.Wait()
	inputPattern := dir+"/image%01d.png"
	images.CreateVideo(inputPattern, "/Users/josephcocozza/Downloads/OUTPUT.mp4",10)
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

