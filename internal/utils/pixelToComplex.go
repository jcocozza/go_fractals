package utils

import "image"

// convert a point in the image to a complex point
func PixelToComplex(point *image.Point, width, height int, zoom float64, center complex128) complex128 {

	x := point.X
	y := point.Y

	z := complex(
		float64(x-width/2)/float64(width)*zoom+real(center),
		float64(y-height/2)/float64(height)*zoom+imag(center),
	)
	return z
}

func PointListToComplexList(lst []*image.Point, width, height int, zoom float64, center complex128) []complex128 {
	cList := []complex128{}
	for _, pt := range lst {
		cList = append(cList, PixelToComplex(pt, width, height, zoom, center))
	}
	return cList
}