package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jcocozza/go_fractals/cmd"
	"github.com/jcocozza/go_fractals/internal/utils"
)

type Mandelbrot struct {
    Equation string `json:"equation"`
    MaxEscapeIterations int `json:"max_escape_iterations"`
    Zoom float64 `json:"zoom"`
    CenterPoint ComplexFromTS `json:"center_point"`
    ColoringAlgorithm string `json:"coloring_algorithm"`
    Image bool	`json:"image"`
    Video bool `json:"video"`
    Stl bool `json:"stl"`
    InitialPoint ComplexFromTS `json:"initial_point"`
    ComplexIncrement ComplexFromTS `json:"increment"`
    NumberIncrements int `json:"num_increments"`
    Fps int `json:"fps"`
    Width int `json:"width"`
    Height int `json:"height"`
	FilePath string `json:"file_path"`
}

func MandelbrotControlFlow(m Mandelbrot, uuid string) ([]string, []string, []string) {
	cmdArgsImg := []string{}
	cmdArgsVid := []string{}
	cmdArgsStl := []string{}

	if m.Image {
		cmdArgsImg = []string{}
		cmdArgsImg = append(cmdArgsImg, "mandelbrot")
		cmdArgsImg = append(cmdArgsImg, []string{
			"--centerPoint", m.CenterPoint.ToString(),
			"--equation", m.Equation,
			"--maxItr", fmt.Sprint(m.MaxEscapeIterations),
			"--zoom", fmt.Sprint(m.Zoom),}...
		)
		//if coloredEntry.Checked {
		//	cmdArgs = append(cmdArgs, "--color")
		//}
		cmdArgsImg = append(cmdArgsImg, []string{
			"--height", fmt.Sprint(m.Height),
			"--width", fmt.Sprint(m.Width),}...
		)
	}

	if m.Video {

	}

	if m.Stl {

	}

	cmdArgsImg = append(cmdArgsImg, "--filePath", m.FilePath + "/" + uuid + ".png")
	cmdArgsVid = append(cmdArgsVid, "--filePath", m.FilePath + "/" + uuid + ".mp4")
	cmdArgsStl = append(cmdArgsStl, "--filePath", m.FilePath + "/" + uuid + ".stl")

	return cmdArgsImg, cmdArgsVid, cmdArgsStl
}

func MandelbrotHandler(ctx *gin.Context) {
	var mandelbrotRequest Mandelbrot
	if err := ctx.ShouldBindJSON(&mandelbrotRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uuid := utils.GenerateUUID()
	image, video, stl := MandelbrotControlFlow(mandelbrotRequest, uuid)
	if len(image) > 2 {
		_, err := cmd.RunCmdArgs(image)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if len(video) > 2 {
		_, err := cmd.RunCmdArgs(video)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if len(stl) > 2 {
		_, err := cmd.RunCmdArgs(stl)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, uuid)
}