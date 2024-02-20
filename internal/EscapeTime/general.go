package EscapeTime

import (
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// return true if the input has escaped
type escapeCondition func(complex128) bool

// return a color based on escape time
type ColorGenerator func(int) color.RGBA

// add text to an image
func DrawText(img *image.RGBA, text string, x, y int) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}
	d.DrawString(text)
}

// return grey scale based on number of iterations
func GreyScale(itr int) color.RGBA {
	grayColor := color.Gray{Y: uint8(itr % 256)}
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

func BurningColor(itr int) color.RGBA {
	return color.RGBA{
		R: 255,
		G: uint8((itr * 7) % 256),
		B: 0,
		A: uint8((itr * 15) % 256),
	}
}

func InfernoColor(itr int) color.RGBA {
	const scale float64 = 0.2 // Adjust this scale for a different appearance

	// Map the iteration value to the Inferno color map
	r := uint8(255 * scale)
	g := uint8(((float64(itr) * 5) / 256) * 256 * scale)
	b := uint8(((float64(itr) * 10) / 256) * 256 * scale)

	return color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}
}