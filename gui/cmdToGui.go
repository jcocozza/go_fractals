package gui

import (
	"log/slog"

	"github.com/jcocozza/go_fractals/cmd"
	"github.com/jcocozza/go_fractals/internal/utils"
)

func ifsControlFlow(cmdArgs []string, outputType, algorithm string) []string {
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
	return cmdArgs
}

func juliaControlFlow(cmdArgs []string, outputType, algorithm string) []string {
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
	return cmdArgs
}

func mandelbrotControlFlow(cmdArgs []string) []string {
	cmdArgs = append(cmdArgs, "mandelbrot")
	cmdArgs = append(cmdArgs, []string{
					"--centerPoint", getEntryValue(centerPointEntry),
					"--equation", getEntryValue(mandelbrotEqnEntry),
					"--maxItr", getEntryValue(maxItrEntry),
					"--zoom", getEntryValue(zoomEntry),}...
					)
	return cmdArgs
}

func generalControlFlow(cmdArgs []string) []string {
	if coloredEntry.Checked {
		cmdArgs = append(cmdArgs, "--color")
	}
	cmdArgs = append(cmdArgs, []string{
		"--fileName", getEntryValue(fileNameEntry),
		"--height", getEntryValue(heightEntry),
		"--width", getEntryValue(widthEntry),}...
	)
	return cmdArgs
}

func runCmdArgs(cmdArgs []string) {
	slog.Info("Running command: " + utils.StrListToString(cmdArgs))
	cmd.RootCmd.SetArgs(cmdArgs)
	if err := cmd.RootCmd.Execute(); err != nil {
		slog.Error("error", err)
	}
}

func cmdOnGui(fractalType, outputType, algorithm string) {
	slog.Info("generating fractal")
	var cmdArgs []string
	if fractalType == "ifs" {
		cmdArgs = ifsControlFlow(cmdArgs, outputType, algorithm)
	} else if fractalType == "julia" {
		cmdArgs = juliaControlFlow(cmdArgs, outputType, algorithm)
	} else if fractalType == "mandelbrot" {
		cmdArgs =  mandelbrotControlFlow(cmdArgs)
	}
	cmdArgs = generalControlFlow(cmdArgs)
	runCmdArgs(cmdArgs)
}