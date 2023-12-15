package gui

import (
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/jcocozza/go_fractals/utils"
)

var (
	startOverBtn *widget.Button
)

func showPage(w fyne.Window, content fyne.CanvasObject) {
    w.SetContent(container.New(layout.NewVBoxLayout(), content))
}

func initStartOverBtn(w fyne.Window, startPage *fyne.Container) {
	startOverBtn = widget.NewButton("Start Over", func() {
		showPage(w, startPage)
	})
}

func page1(w fyne.Window) *fyne.Container {
	fractalTypes := []string{"ifs", "julia", "mandelbrot"}
	fractalChoices = widget.NewRadioGroup(fractalTypes, func(value string) {
		fractalType = value
		slog.Info("Set fractal choice to: " + fractalType)
	})
	return container.NewVBox(
		widget.NewLabel("Choose fractal type"),
		fractalChoices,
	)
}

func page2(w fyne.Window, fractalType string) *fyne.Container {
	outputTypes := []string{"image","video","stl"}
	outputChoices := widget.NewRadioGroup(outputTypes, func(s string) {
		outputType = s
		slog.Info("Set output type to: " + outputType)
	})

	return container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Fractal type chosen: %s", fractalType)),
		widget.NewLabel("Choose output type"),
		outputChoices,
	)
}

func page3formControl(pg *fyne.Container, fractalType, outputType string) {
	slog.Info("fractalType: " + fractalType)
	slog.Info("outputType: " + outputType)
	if fractalType == "ifs" {
		pg.Add(ifsForm)
		if outputType == "image" {

		} else if outputType == "video" {
			pg.Add(videoForm)

		} else if outputType == "stl" {
			pg.Add(stlForm)
		}
	} else if fractalType == "julia" {
		pg.Add(juliaForm)
		if outputType == "image" {

		} else if outputType == "video" {
			pg.Add(juliaEvolveForm)
			pg.Add(videoForm)
		} else if outputType == "stl" {
			pg.Add(juliaEvolveForm)
			pg.Add(stlForm)
		}
	} else if fractalType == "mandelbrot" {
		pg.Add(mandelbrotForm)
	}
}

func page3(w fyne.Window, fractalType, outputType string) *fyne.Container {
	pg := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Fractal type chosen: %s. Output Type chosen: %s", fractalType, outputType)),
		generalForm,
	)
	page3formControl(pg, fractalType, outputType)


	generateButton := widget.NewButton("Generate Fractal", func() {
		slog.Info("generating fractal")
		cmdOnGui(fractalType, outputType, algorithm)
		showPage(w, page4())
	})

	pg.Add(generateButton)
	return pg
}

func page4() *fyne.Container {
	pg := container.NewVBox(
		handleOutPut(fileNameEntry, outputType),
		startOverBtn,
	)
	return pg
}

func handleOutPut(fileNameEntry *widget.Entry, outputType string) *fyne.Container {

	var containerOut *fyne.Container

	filePath := utils.GetDownloadDir() + "/" + getEntryValue(fileNameEntry)
	slog.Info("loading fractal image from: "+ filePath)
	if outputType == "image" {
		filePath += ".png"
		img := canvas.NewImageFromFile(filePath)
		img.FillMode = canvas.ImageFillOriginal

		containerOut = container.NewVBox(
			img,
			widget.NewButton("Open File", func() {
				utils.Open(filePath)
			}),
		)
	} else if outputType == "video" {
		filePath += ".mp4"
		containerOut = container.NewVBox(
			widget.NewCard("Video Saved to: " + filePath, "", nil),
			widget.NewButton("Open File", func() {
				utils.Open(filePath)
			}),
		)
	} else if outputType == "stl" {
		filePath += ".stl"
		containerOut = container.NewVBox(
			widget.NewCard("stl Saved to: " + filePath, "", nil),
			widget.NewButton("Open File", func() {
				utils.Open(filePath)
			}),
		)
	}

	return containerOut
}