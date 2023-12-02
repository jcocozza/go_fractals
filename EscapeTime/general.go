package EscapeTime

import (
	"image/color"
)

const (
	width = 200
	height = 200
	maxIterations = 1000
)

// return true if the input has escaped
type escapeCondition func(complex128) bool

// return a color based on escape time
type colorGenerator func(int) color.RGBA

// return grey scale based on number of iterations
func GreyScale(itr int) color.RGBA {
	grayColor := color.Gray{Y: uint8(itr % 256)}

	// Check if the grayscale color is black
	//if grayColor.Y == 0 {
	//	return color.RGBA{0, 0, 0, 0} // Transparent color
	//}

	return color.RGBA{grayColor.Y, grayColor.Y, grayColor.Y, 255}
}

// return grey scale based on number of iterations
func GreyScaleClear(itr int) color.RGBA {
	grayColor := color.Gray{Y: uint8(itr % 256)}

	// Check if the grayscale color is black
	if grayColor.Y == 0 {
		return color.RGBA{0, 0, 0, 0} // Transparent color
	}

	return color.RGBA{grayColor.Y, grayColor.Y, grayColor.Y, 255}
}

func UnitaryScale(itr int) color.RGBA {
	return color.RGBA{0,255,0, 255} //green
}

// return a color based on number of iterations
func GenerateColor(itr int) color.RGBA {
	return color.RGBA{
		R: uint8((itr * 37) % 256),
		G: uint8((itr * 73) % 256),
		B: uint8((itr * 139) % 256),
		A: 255,
	}
}