package primitives

import (
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "os"
)

func Plot() {
    // Define the width and height of the image
    width := 400
    height := 400

    // Create a new RGBA image
    img := image.NewRGBA(image.Rect(0, 0, width, height))

    // Choose a color for the parallelogram (e.g., green)
    parallelogramColor := color.RGBA{0, 255, 0, 255} // Green color

    // Define the coordinates of the parallelogram's vertices
    vertices := []image.Point{
        image.Point{0, 0},
        image.Point{100, 0},
        image.Point{50, 200},
        image.Point{150, 200},
    }

    // Fill the parallelogram with the specified color
    drawPolygon(img, vertices, parallelogramColor)

    // Create a PNG file to save the image
    file, err := os.Create("colored_parallelogram.png")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // Encode and save the image as a PNG file
    err = png.Encode(file, img)
    if err != nil {
        panic(err)
    }
}

// drawPolygon fills a polygon with the given color using the scanline algorithm.
func drawPolygon(img *image.RGBA, vertices []image.Point, color color.Color) {
    // Find the bounding box of the polygon
    min, max := findBoundingBox(vertices)

    // Iterate through the bounding box and fill the polygon
    for y := min.Y; y <= max.Y; y++ {
        intersections := findIntersections(vertices, y)
        for i := 0; i < len(intersections); i += 2 {
            x0 := clamp(intersections[i], min.X, max.X)
            x1 := clamp(intersections[i+1], min.X, max.X)
            if x0 <= x1 {
                draw.Draw(img, image.Rect(x0, y, x1+1, y+1), &image.Uniform{color}, image.Point{}, draw.Src)
            }
        }
    }
}

// findBoundingBox finds the minimum and maximum coordinates of the polygon's bounding box.
func findBoundingBox(vertices []image.Point) (min, max image.Point) {
    min = vertices[0]
    max = vertices[0]
    for _, v := range vertices {
        if v.X < min.X {
            min.X = v.X
        }
        if v.X > max.X {
            max.X = v.X
        }
        if v.Y < min.Y {
            min.Y = v.Y
        }
        if v.Y > max.Y {
            max.Y = v.Y
        }
    }
    return min, max
}

// findIntersections finds the x-coordinates where a horizontal line at the specified y value intersects the edges of the polygon.
func findIntersections(vertices []image.Point, y int) []int {
    var intersections []int
    for i, j := 0, len(vertices)-1; i < len(vertices); i++ {
        vi := vertices[i]
        vj := vertices[j]
        if vi.Y < vj.Y {
            vi, vj = vj, vi
        }
        if vi.Y >= y && vj.Y < y {
            x := vi.X - int(float64(vi.Y-y)*float64(vi.X-vj.X)/float64(vi.Y-vj.Y))
            intersections = append(intersections, x)
        }
        j = i
    }
    return intersections
}

// clamp clamps the value to the specified range.
func clamp(value, min, max int) int {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}
