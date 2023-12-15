package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	ifsForm *widget.Form
	mandelbrotForm *widget.Form
	juliaForm *widget.Form
	juliaEvolveForm *widget.Form
	generalForm *widget.Form

	stlForm *fyne.Container
	videoForm *fyne.Container
)

func initForms() {
	ifsForm = &widget.Form{
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
	mandelbrotForm = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Equation", Widget: mandelbrotEqnEntry},
			{Text: "zoom", Widget: zoomEntry},
			{Text: "max iterations for escape:", Widget: maxItrEntry},
			{Text: "complex center point:", Widget: centerPointEntry},
			{Text: "", Widget: coloredEntry},},
		OnSubmit: func() {},
	}
	juliaForm = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Equation", Widget: juliaEqnEntry},
			{Text: "zoom", Widget: zoomEntry},
			{Text: "max iterations for escape:", Widget: maxItrEntry},
			{Text: "complex center point:", Widget: centerPointEntry},
			{Text: "", Widget: coloredEntry},},
		OnSubmit: func() {},
	}
	juliaEvolveForm = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "inital complex point", Widget: cInitEntry},
			{Text: "complex incremenent", Widget: cIncrementEntry},
			{Text: "number increments", Widget: numIncrementsEntry},
		},
		OnSubmit: func() {},
	}
	generalForm = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "width", Widget: widthEntry},
			{Text: "height", Widget: heightEntry},
			{Text: "file name", Widget: fileNameEntry},
		},
		OnSubmit: func() {},
	}

	stlForm = container.NewHBox(threeDimensionalEntry, solidEntry, binaryEntry,)
	videoForm = container.NewHBox(widget.NewLabel("fps:"), fpsEntry,)
}