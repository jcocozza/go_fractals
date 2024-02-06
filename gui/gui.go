package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func initPages(w fyne.Window) (*fyne.Container,*fyne.Container,*fyne.Container,*fyne.Container) {
	p3 := container.NewVBox()
	p2 := container.NewVBox(
		page2(w, fractalType),
		widget.NewButton("Continue", func() {
			p3.Add(page3(w, fractalType, outputType),)
			p3.Add(startOverBtn)
            showPage(w, p3)
        }),
	)

	p1 := container.NewVBox(
		page1(w),
		widget.NewButton("Continue", func() {
            showPage(w, p2)
        }),
	)

	p0 := container.NewVBox(
		widget.NewCard("Welcome to go fractals", "",
		widget.NewButton("start", func() {
			showPage(w, p1)
		}),
	),
	widget.NewButton("Draw Mandelbrot evolution", func ()  {
		cont := DrawOnImage(w)
		showPage(w, cont)
	}),
)

	// TODO ensure that the start over button actually resets things
	startOverBtn = widget.NewButton("Start Over", func() {
		showPage(w, p0)
		initPages(w)
	})

	//p0.Add(complexValueSlider(-10,10,0.1))
	p1.Add(startOverBtn)
	p2.Add(startOverBtn)

	return p0,p1,p2,p3
}

func GUI() {
	a := app.New()
	w := a.NewWindow("Go Fractals")

	initWidgets()
	initForms()

	p0,_,_,_ := initPages(w)
	showPage(w, p0)

	//w.SetContent(cont)
	w.ShowAndRun()
}