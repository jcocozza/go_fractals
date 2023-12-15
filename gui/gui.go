package gui

import (
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"github.com/jcocozza/go_fractals/cmd"
	"github.com/jcocozza/go_fractals/utils"
)

func showPage(w fyne.Window, content fyne.CanvasObject) {
    w.SetContent(container.New(layout.NewVBoxLayout(), content))
}

func getEntryValue(w *widget.Entry) string {
	if w.Text == "" {
		return w.PlaceHolder
	}
	return w.Text
}

func loadImage(fileNameEntry *widget.Entry, outputType string) *fyne.Container {

	var containerOut *fyne.Container

	filePath := utils.GetDownloadDir() + "/" + getEntryValue(fileNameEntry)
	slog.Info("loading fractal image from: "+ filePath)
	if outputType == "image" {
		filePath += ".png"
		img := canvas.NewImageFromFile(filePath)
		img.FillMode = canvas.ImageFillOriginal

		containerOut = container.NewHBox(
			img,
		)
	} else if outputType == "video" {
		filePath += ".mp4"
		containerOut = container.NewHBox(
			widget.NewCard("Video Saved to: " + filePath, "", nil),
		)
	} else if outputType == "stl" {
		filePath += ".stl"
		containerOut = container.NewHBox(
			widget.NewCard("stl Saved to: " + filePath, "", nil),
		)
	}

	return containerOut
}

func GUI() {
	a := app.New()
	w := a.NewWindow("Go Fractals")

	// general
	fpsEntry := widget.NewEntry()
	fpsEntry.SetPlaceHolder(fmt.Sprintf("%d",cmd.FpsDefault))

	fileNameEntry := widget.NewEntry()
	fileNameEntry.SetPlaceHolder(cmd.FileNameDefault)

	widthEntry := widget.NewEntry()
	widthEntry.SetPlaceHolder(fmt.Sprintf("%d",cmd.WidthDefault))

	heightEntry := widget.NewEntry()
	heightEntry.SetPlaceHolder(fmt.Sprintf("%d",cmd.HeightDefault))

	// ifs
	ifsPathEntry := widget.NewEntry()
	ifsPathEntry.SetPlaceHolder(cmd.IfsPathDefault)

	var algorithm string
	algorithms := []string{"probabilistic", "deterministic"}
	algorithmChoices := widget.NewRadioGroup(algorithms, func(s string) {
		algorithm = s
	})

	numItrsEntry := widget.NewEntry()
	numItrsEntry.SetPlaceHolder(string(cmd.NumIterationsDefault))

	numPointsEntry := widget.NewEntry()
	numPointsEntry.SetPlaceHolder(string(cmd.NumPointsDefault))

	probabilitiesEntry := widget.NewEntry()
	probabilitiesEntry.SetPlaceHolder(utils.ListToString(cmd.ProbabilitiesListDefault))

	numStacksEntry := widget.NewEntry()
	numStacksEntry.SetPlaceHolder(string(cmd.NumStacksDefault))

	thicknessEntry := widget.NewEntry()
	thicknessEntry.SetPlaceHolder("15")

	// julia & mandelbrot
	juliaEqnEntry := widget.NewEntry()
	juliaEqnEntry.SetPlaceHolder(cmd.JuliaEquationDefault)

	mandelbrotEqnEntry := widget.NewEntry()
	mandelbrotEqnEntry.SetPlaceHolder(cmd.MandelbrotEquationDefault)

	coloredEntry := widget.NewCheck("colored?", nil)
	cInitEntry := widget.NewEntry()
	cInitEntry.SetPlaceHolder(cmd.CINitStringDefault)

	cIncrementEntry := widget.NewEntry()
	cIncrementEntry.SetPlaceHolder(cmd.CIncrementStringDefault)

	numIncrementsEntry := widget.NewEntry()
	numIncrementsEntry.SetPlaceHolder(fmt.Sprintf("%d", cmd.NumIncrementsDefault))

	// stl stuff
	threeDimensionalEntry := widget.NewCheck("generate stl?", nil)
	solidEntry := widget.NewCheck("generate stl as solid?", nil)
	binaryEntry := widget.NewCheck("generate stl as binary?", nil)

	centerPointEntry := widget.NewEntry()
	centerPointEntry.SetPlaceHolder(cmd.CenterPointStringDefault)

	zoomEntry := widget.NewEntry()
	zoomEntry.SetPlaceHolder(fmt.Sprintf("%f", cmd.ZoomDefault))

	maxItrEntry := widget.NewEntry()
	maxItrEntry.SetPlaceHolder(fmt.Sprintf("%d", cmd.MaxItrDefault))

	stlForm := container.NewHBox(threeDimensionalEntry, solidEntry, binaryEntry,)
	videoForm := container.NewHBox(widget.NewLabel("fps:"), fpsEntry,)

	var page1 *fyne.Container
	var page2 func(string) *fyne.Container
	var page3 func(string,string) *fyne.Container

	ifsForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ifs path", Widget: ifsPathEntry},
			{Text: "algorithm", Widget: algorithmChoices},
			{Text: "number of iterations", Widget: numItrsEntry},
			{Text: "number of initial points", Widget: numPointsEntry},
			{Text: "probability list", Widget: probabilitiesEntry},
			{Text: "number of stacks", Widget: numStacksEntry},
			{Text: "thickness", Widget: thicknessEntry},
		},
		OnSubmit: func() {},
	}
	mandelbrotForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Equation", Widget: mandelbrotEqnEntry},
			{Text: "zoom", Widget: zoomEntry},
			{Text: "max iterations for escape:", Widget: maxItrEntry},
			{Text: "complex center point:", Widget: centerPointEntry},
			{Text: "", Widget: coloredEntry},},
		OnSubmit: func() {},
	}
	juliaForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Equation", Widget: juliaEqnEntry},
			{Text: "zoom", Widget: zoomEntry},
			{Text: "max iterations for escape:", Widget: maxItrEntry},
			{Text: "complex center point:", Widget: centerPointEntry},
			{Text: "", Widget: coloredEntry},},
		OnSubmit: func() {},
	}
	juliaEvolveForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "inital complex point", Widget: cInitEntry},
			{Text: "complex incremenent", Widget: cIncrementEntry},
			{Text: "number increments", Widget: numIncrementsEntry},
		},
		OnSubmit: func() {},
	}
	generalForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "width", Widget: widthEntry},
			{Text: "height", Widget: heightEntry},
			{Text: "file name", Widget: fileNameEntry},
		},
		OnSubmit: func() {},
	}

	var fractalType string
	fractalTypes := []string{"ifs", "julia", "mandelbrot"}
	fractalChoices := widget.NewRadioGroup(fractalTypes, func(value string) {
        fractalType = value
    })
	page1 = container.NewVBox(
        widget.NewLabel("Choose fractal type"),
        fractalChoices,
        widget.NewButton("Continue", func() {
            showPage(w, page2(fractalType))
        }),
    )

	startOverBtn := widget.NewButton("Start Over", func() {
		showPage(w, page1)
	})

	var outputType string
	outputTypes := []string{"image","video","stl"}
	outputChoices := widget.NewRadioGroup(outputTypes, func(s string) {
		outputType = s
	})
	page2 = func(fractalType string) *fyne.Container {
        return container.NewVBox(
            widget.NewLabel(fmt.Sprintf("Fractal type chosen: %s", fractalType)),
            widget.NewLabel("Choose output type"),
            outputChoices,
            widget.NewButton("Go to Page 3", func() {
                showPage(w, page3(fractalType, outputType))
            }),
			startOverBtn,
        )
    }

	page3 = func(fractalType, outputType string) *fyne.Container {
		pg3 := container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Fractal type chosen: %s. Output Type chosen: %s", fractalType, outputType)),
			generalForm,
		)
		if fractalType == "ifs" {
			pg3.Add(ifsForm)
			if outputType == "image" {

			} else if outputType == "video" {
				pg3.Add(videoForm)

			} else if outputType == "stl" {
				pg3.Add(stlForm)
			}
		} else if fractalType == "julia" {
			pg3.Add(juliaForm)
			if outputType == "image" {

			} else if outputType == "video" {
				pg3.Add(juliaEvolveForm)
				pg3.Add(videoForm)
			} else if outputType == "stl" {
				pg3.Add(juliaEvolveForm)
				pg3.Add(stlForm)
			}
		} else if fractalType == "mandelbrot" {
			pg3.Add(mandelbrotForm)
		}

		generateButton := widget.NewButton("Generate Fractal", func() {
			slog.Info("generating fractal")
			var cmdArgs []string
			if fractalType == "ifs" {
				cmdArgs = append(cmdArgs, "ifs")
				if outputType == "image" {
					cmdArgs = append(cmdArgs, "image")
				} else if outputType == "video" {
					cmdArgs = append(cmdArgs, "evolve")
				} else if outputType == "stl" {
					cmdArgs = append(cmdArgs, "evolve")

					cmdArgs = append(cmdArgs, []string{
						"--numStacks", getEntryValue(numStacksEntry),
						"--thickness", getEntryValue(thicknessEntry),
					}...)
				}

				if algorithm == "probabilistic" {
					cmdArgs = append(cmdArgs, "--algo-p")
				} else if algorithm == "deterministic" {
					cmdArgs = append(cmdArgs, "--algo-d")
				}
				cmdArgs = append(cmdArgs, []string{
					"--path", getEntryValue(ifsPathEntry),
					"--numItr", getEntryValue(numItrsEntry),
					"--numPoints", getEntryValue(numPointsEntry),
				}...)

				if getEntryValue(probabilitiesEntry) != "" {
					cmdArgs = append(cmdArgs,
						[]string{"--probabilities", getEntryValue(probabilitiesEntry)}...,)
				}
			} else if fractalType == "julia" {
				if outputType == "image" {
					cmdArgs = append(cmdArgs, "julia")

					cmdArgs = append(cmdArgs, []string{
						"--centerPoint", getEntryValue(centerPointEntry),
						"--equation", getEntryValue(juliaEqnEntry),
						"--maxItr", getEntryValue(maxItrEntry),
						"--zoom", getEntryValue(zoomEntry),}...
					)

				} else if outputType == "video" {
					cmdArgs = append(cmdArgs, "julia-evolve")
					cmdArgs = append(cmdArgs, []string{
						"--centerPoint", getEntryValue(centerPointEntry),
						"--complexIncrement", getEntryValue(cIncrementEntry),
						"--equation", getEntryValue(mandelbrotEqnEntry),
						"--fps", getEntryValue(fpsEntry),
						"--initialComplex", getEntryValue(cInitEntry),
						"--maxItr", getEntryValue(maxItrEntry),
						"--numIncrements", getEntryValue(numIncrementsEntry),
						"--zoom", getEntryValue(zoomEntry),}...
					)

				} else if outputType == "stl" {
					cmdArgs = append(cmdArgs, "julia-evolve")
					cmdArgs = append(cmdArgs, []string{
						"--centerPoint", getEntryValue(centerPointEntry),
						"--complexIncrement", getEntryValue(cIncrementEntry),
						"--equation", getEntryValue(mandelbrotEqnEntry),
						"--fps", getEntryValue(fpsEntry),
						"--initialComplex", getEntryValue(cInitEntry),
						"--maxItr", getEntryValue(maxItrEntry),
						"--numIncrements", getEntryValue(numIncrementsEntry),
						"--zoom", getEntryValue(zoomEntry),}...
					)

					if binaryEntry.Checked {
						cmdArgs = append(cmdArgs, "--binary")
					}
					if solidEntry.Checked {
						cmdArgs = append(cmdArgs, "--solid")
					}
					cmdArgs = append(cmdArgs, "--threeDim")
				}
			} else if fractalType == "mandelbrot" {
				cmdArgs = append(cmdArgs, "mandelbrot")
				cmdArgs = append(cmdArgs, []string{
					"--centerPoint", getEntryValue(centerPointEntry),
					"--equation", getEntryValue(mandelbrotEqnEntry),
					"--maxItr", getEntryValue(maxItrEntry),
					"--zoom", getEntryValue(zoomEntry),}...
				)
			}

			if coloredEntry.Checked {
				cmdArgs = append(cmdArgs, "--color")
			}
			cmdArgs = append(cmdArgs, []string{
				"--fileName", getEntryValue(fileNameEntry),
				"--height", getEntryValue(heightEntry),
				"--width", getEntryValue(widthEntry),}...
			)
			slog.Info("Running command: ", cmdArgs)

			cmd.RootCmd.SetArgs(cmdArgs)
			if err := cmd.RootCmd.Execute(); err != nil {
				slog.Error("error", err)
			}
			page4 := container.NewVBox(
				loadImage(fileNameEntry, outputType),
				startOverBtn,
			)
			showPage(w, page4)
		})

		pg3.Add(generateButton)
		pg3.Add(startOverBtn)
		return pg3
	}

	showPage(w, page1)
	w.ShowAndRun()
}