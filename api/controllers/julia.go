package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jcocozza/go_fractals/cmd"
	"github.com/jcocozza/go_fractals/internal/utils"
)

type Julia struct {
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

func JuliaControlFlow(j Julia, uuid string) ([]string, []string, []string) {
	cmdArgsImg := []string{}
	cmdArgsVid := []string{}
	cmdArgsStl := []string{}

	if j.Image {
		cmdArgsImg = []string{}
		cmdArgsImg = append(cmdArgsImg, "julia")
		cmdArgsImg = append(cmdArgsImg, []string{
			"--centerPoint", j.CenterPoint.ToString(),
			"--equation", j.Equation,
			"--maxItr", fmt.Sprint(j.MaxEscapeIterations),
			"--zoom", fmt.Sprint(j.Zoom),}...
		)

		//if coloredEntry.Checked {
		//	cmdArgs = append(cmdArgs, "--color")
		//}
		cmdArgsImg = append(cmdArgsImg, []string{
			"--height", fmt.Sprint(j.Height),
			"--width", fmt.Sprint(j.Width),}...
		)
	}

	if j.Video {
		cmdArgsVid = []string{}
		cmdArgsVid = append(cmdArgsVid, "julia-evolve")
		cmdArgsVid = append(cmdArgsVid, []string{
			"--centerPoint", j.CenterPoint.ToString(),
			"--complexIncrement", j.ComplexIncrement.ToString(),
			"--equation", j.Equation,
			"--fps", fmt.Sprint(j.Fps),
			"--initialComplex", j.InitialPoint.ToString(),
			"--maxItr", fmt.Sprint(j.MaxEscapeIterations),
			"--numIncrements", fmt.Sprint(j.NumberIncrements),
			"--zoom", fmt.Sprint(j.Zoom),}...
		)

		cmdArgsVid = append(cmdArgsVid, []string{
			"--height", fmt.Sprint(j.Height),
			"--width", fmt.Sprint(j.Width),}...
		)
	}
	if j.Stl {
		cmdArgsStl = []string{}
		cmdArgsStl = append(cmdArgsStl, "julia-evolve")
		cmdArgsStl = append(cmdArgsStl, []string{
			"--centerPoint", j.CenterPoint.ToString(),
			"--complexIncrement", j.ComplexIncrement.ToString(),
			"--equation", j.Equation,
			"--initialComplex", j.InitialPoint.ToString(),
			"--maxItr", fmt.Sprint(j.MaxEscapeIterations),
			"--numIncrements", fmt.Sprint(j.NumberIncrements),
			"--zoom",  fmt.Sprint(j.Zoom),}...
		)
		cmdArgsStl = append(cmdArgsStl, "--threeDim")
	}

	cmdArgsImg = append(cmdArgsImg, "--filePath", j.FilePath + "/" + uuid + ".png")
	cmdArgsVid = append(cmdArgsVid, "--filePath", j.FilePath + "/" + uuid + ".mp4")
	cmdArgsStl = append(cmdArgsStl, "--filePath", j.FilePath + "/" + uuid + ".stl")

	return cmdArgsImg, cmdArgsVid, cmdArgsStl
}

func JuliaHandler(ctx *gin.Context) {
	var juliaRequest Julia
	if err := ctx.ShouldBindJSON(&juliaRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uuid := utils.GenerateUUID()
	image, video, stl := JuliaControlFlow(juliaRequest, uuid)

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