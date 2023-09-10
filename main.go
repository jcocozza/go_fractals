package main

import (
	"gonum.org/v1/gonum/mat"
	IFS "github.com/jcocozza/go_fractals/IteratedFunctionSystems"
	viz "github.com/jcocozza/go_fractals/visualizer"
)
func main() {
	// Square Fractal Thingy Transformation
	T1, _ := IFS.NewTransformation(*mat.NewDense(2,2, []float64{.5,0,0,.5}), *mat.NewDense(2,1, []float64{0,0}))
	T2, _ := IFS.NewTransformation(*mat.NewDense(2,2, []float64{-.5,0,0,.5}), *mat.NewDense(2,1, []float64{1,0}))
	T3, _ := IFS.NewTransformation(*mat.NewDense(2,2, []float64{.5,0,0,-.5}), *mat.NewDense(2,1, []float64{0,1}))
	T4, err := IFS.NewTransformation(*mat.NewDense(2,2, []float64{.25,0,0,.25}), *mat.NewDense(2,1, []float64{.75,.75}))

	/*
	T1, _ := IFS.NewTransformation(*mat.NewDense(2,2, []float64{0,0,0,.16}), *mat.NewDense(2,1, []float64{0,0}))
	T2, _ := IFS.NewTransformation(*mat.NewDense(2,2, []float64{.85,0.04,-.04,.85}), *mat.NewDense(2,1, []float64{0,1.6}))
	T3, _ := IFS.NewTransformation(*mat.NewDense(2,2, []float64{0.2,-.26,.23,.22}), *mat.NewDense(2,1, []float64{0,1.6}))
	T4, err := IFS.NewTransformation(*mat.NewDense(2,2, []float64{-.15,.28,.26,.24}), *mat.NewDense(2,1, []float64{0,.44}))
	*/
	/*
	if err == nil {
		IFS := IFS.NewIteratedFunctionSystem([]IFS.Transformation{*T1,*T2,*T3,*T4}, 9, 2)
		pointsList := IFS.RunDeterministic()

		const width, height = 1000,1000
		fractal := viz.NewFractalImage(width, height, "my_fractal.png", pointsList)
		fractal.WriteImage()


	} else {
		print("Error:", err)
	}
	*/

	if err == nil {
		ifsys := IFS.NewIteratedFunctionSystem([]IFS.Transformation{*T1,*T2,*T3,*T4}, 9, 2)

		const width, height = 1000,1000
		fv := viz.NewFractalVideo(width, height, "my_fractal_video.mp4", *ifsys)

		fv.WriteVideoImages(10)
		fv.CreateVideo()

	} else {
		print("Error:", err)
	}
}

