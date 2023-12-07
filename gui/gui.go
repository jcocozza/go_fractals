package gui

import (
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jcocozza/go_fractals/cmd"
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

		//cmdArgs = []string{"julia","-e", "1/(z*z + .72i)", "-F", "testGUI"}
		generateButton := widget.NewButton("Generate Fractal", func() {
			slog.Info("generating fractal")
			var cmdArgs []string
			if fractalType == "ifs" {
				cmdArgs = append(cmdArgs, "ifs")
				if outputType == "image" {
				} else if outputType == "video" {
				} else if outputType == "stl" {
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
						"--equation", getEntryValue(juliaEqnEntry),
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
						"--equation", getEntryValue(juliaEqnEntry),
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
		})

		pg3.Add(generateButton)
		pg3.Add(startOverBtn)

		return pg3
	}


	showPage(w, page1)
	w.ShowAndRun()
}