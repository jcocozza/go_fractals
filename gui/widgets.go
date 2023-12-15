package gui

import (
	"log/slog"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/jcocozza/go_fractals/cmd"
	"github.com/jcocozza/go_fractals/utils"
)

var (
	fractalChoices *widget.RadioGroup
	fpsEntry *widget.Entry
	fileNameEntry *widget.Entry
	widthEntry *widget.Entry
	heightEntry *widget.Entry
	ifsPathEntry *widget.Entry
	algorithmChoices *widget.RadioGroup
	numItrsEntry *widget.Entry
	numPointsEntry *widget.Entry
	probabilitiesEntry *widget.Entry
	numStacksEntry *widget.Entry
	thicknessEntry *widget.Entry
	juliaEqnEntry *widget.Entry
	mandelbrotEqnEntry *widget.Entry
	coloredEntry *widget.Check
	cInitEntry *widget.Entry
	cIncrementEntry *widget.Entry
	numIncrementsEntry *widget.Entry
	threeDimensionalEntry *widget.Check
	solidEntry *widget.Check
	binaryEntry *widget.Check
	centerPointEntry *widget.Entry
	zoomEntry *widget.Entry
	maxItrEntry *widget.Entry

	algorithm string
	fractalType string
	outputType string
)

func getEntryValue(w *widget.Entry) string {
	if w.Text == "" {
		return w.PlaceHolder
	}
	return w.Text
}

func complexValueSlider(min, max, increment float64) *fyne.Container {
	realSlider := widget.NewSlider(min,max)
	realSlider.Step = increment
	imaginarySlider := widget.NewSlider(min,max)
	imaginarySlider.Step = increment

	r := binding.BindFloat(&realSlider.Value)
	i := binding.BindFloat(&imaginarySlider.Value)

	cValueLabel := widget.NewLabelWithData(binding.NewSprintf("%.2f + %.2fi",r,i))

	realSlider.OnChanged = func(value float64) {
		slog.Debug("Set real value to: " + fmt.Sprint(value))
		r.Set(value)
	}

	imaginarySlider.OnChanged = func(value float64) {
		slog.Debug("Set imaginary value to: " + fmt.Sprint(value))
		i.Set(value)
	}

	slider := container.NewVBox(
		realSlider,
		imaginarySlider,
		cValueLabel,
	)
	return slider
}


func initWidgets() {

	fractalTypes := []string{"ifs", "julia", "mandelbrot"}
	fractalChoices = widget.NewRadioGroup(fractalTypes, func(value string) {
		fractalType = value
	})

	// general
	fpsEntry = widget.NewEntry()
	fpsEntry.SetPlaceHolder(fmt.Sprintf("%d",cmd.FpsDefault))

	fileNameEntry = widget.NewEntry()
	fileNameEntry.SetPlaceHolder(cmd.FileNameDefault)

	widthEntry = widget.NewEntry()
	widthEntry.SetPlaceHolder(fmt.Sprintf("%d",cmd.WidthDefault))

	heightEntry = widget.NewEntry()
	heightEntry.SetPlaceHolder(fmt.Sprintf("%d",cmd.HeightDefault))

	// ifs
	ifsPathEntry = widget.NewEntry()
	ifsPathEntry.SetPlaceHolder(cmd.IfsPathDefault)

	algorithms := []string{"probabilistic", "deterministic"}
	algorithmChoices = widget.NewRadioGroup(algorithms, func(s string) {
		algorithm = s
	})

	numItrsEntry = widget.NewEntry()
	numItrsEntry.SetPlaceHolder(fmt.Sprint(cmd.NumIterationsDefault))

	numPointsEntry = widget.NewEntry()
	numPointsEntry.SetPlaceHolder(fmt.Sprint(cmd.NumPointsDefault))

	probabilitiesEntry = widget.NewEntry()
	probabilitiesEntry.SetPlaceHolder(utils.ListToString(cmd.ProbabilitiesListDefault))

	numStacksEntry = widget.NewEntry()
	numStacksEntry.SetPlaceHolder(fmt.Sprint(cmd.NumStacksDefault))

	thicknessEntry = widget.NewEntry()
	thicknessEntry.SetPlaceHolder("15")

	// julia & mandelbrot
	juliaEqnEntry = widget.NewEntry()
	juliaEqnEntry.SetPlaceHolder(cmd.JuliaEquationDefault)

	mandelbrotEqnEntry = widget.NewEntry()
	mandelbrotEqnEntry.SetPlaceHolder(cmd.MandelbrotEquationDefault)

	coloredEntry = widget.NewCheck("colored?", nil)
	cInitEntry = widget.NewEntry()
	cInitEntry.SetPlaceHolder(cmd.CInitStringDefault)

	cIncrementEntry = widget.NewEntry()
	cIncrementEntry.SetPlaceHolder(cmd.CIncrementStringDefault)

	numIncrementsEntry = widget.NewEntry()
	numIncrementsEntry.SetPlaceHolder(fmt.Sprintf("%d", cmd.NumIncrementsDefault))

	// stl stuff
	threeDimensionalEntry = widget.NewCheck("generate stl?", nil)
	solidEntry = widget.NewCheck("generate stl as solid?", nil)
	binaryEntry = widget.NewCheck("generate stl as binary?", nil)

	centerPointEntry = widget.NewEntry()
	centerPointEntry.SetPlaceHolder(cmd.CenterPointStringDefault)

	zoomEntry = widget.NewEntry()
	zoomEntry.SetPlaceHolder(fmt.Sprintf("%f", cmd.ZoomDefault))

	maxItrEntry = widget.NewEntry()
	maxItrEntry.SetPlaceHolder(fmt.Sprintf("%d", cmd.MaxItrDefault))
}
