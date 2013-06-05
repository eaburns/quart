package geom

// Drawing of 2-dimensional primitives.

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// The Canvas interface encapsulates the functions used to draw
// geometric primitives.  The canvas should be oriented such that
// the lower left corner is the point 0,0, up is positive Y and right is
// positive X.
type Canvas interface {
	Size() (int, int)
	StrokeLine(c color.Color, x0, y0, x1, y1 int)
	FillCircle(c color.Color, x, y, r int)
}

// Draw draws a point on the canvas.
func (pt Point) Draw(cv Canvas, cl color.Color) {
	const radius = 4
	x0, y0 := int(pt[0]+0.5), int(pt[1]+0.5)
	cv.FillCircle(cl, x0, y0, radius)
}

// DrawAt draws the vector extending from a given point.
func (v Vector) DrawAt(cv Canvas, cl color.Color, p Point) {
	p.Draw(cv, cl)
	x0, y0 := int(p[0]+0.5), int(p[1]+0.5)
	p1 := p.Plus(v)
	x1, y1 := int(p1[0]+0.5), int(p1[1]+0.5)
	cv.StrokeLine(cl, x0, y0, x1, y1)
}

// Draw draws a ray on the canvas.
func (ray Ray) Draw(cv Canvas, cl color.Color) {
	const length = 25
	ray.Direction.ScaledBy(length).DrawAt(cv, cl, ray.Origin)
}

// Draw draws a line on the canvas.
func (l Line) Draw(cv Canvas, cl color.Color) {
	wi, hi := cv.Size()
	w, h := float64(wi), float64(hi)
	segs := [4]Line{
		{Origin: Point{0, 0}, Normal: Vector{0, 1}},
		{Origin: Point{0, 0}, Normal: Vector{1, 0}},
		{Origin: Point{w - 1, h - 1}, Normal: Vector{0, -1}},
		{Origin: Point{w - 1, h - 1}, Normal: Vector{-1, 0}},
	}

	var ends []Point
	for _, s := range segs {
		p, hit := l.LineIntersection(s)
		if hit && onCanvas(p, cv) && (len(ends) == 0 || !p.Equals(ends[0])) {
			ends = append(ends, p)
		}
	}

	x0, y0 := int(ends[0][0]+0.5), int(ends[0][1]+0.5)
	x1, y1 := int(ends[1][0]+0.5), int(ends[1][1]+0.5)
	cv.StrokeLine(cl, x0, y0, x1, y1)

	len := ends[0].Distance(ends[1])
	dir := l.Direction()
	p := ends[0].Plus(dir.ScaledBy(len / 2))
	if p[0] < 0 || p[0] >= w || p[1] < 0 || p[1] >= h {
		p = ends[1].Plus(dir.ScaledBy(len / 2))
	}
	Ray{Origin: p, Direction: l.Normal}.Draw(cv, cl)
}

func onCanvas(p Point, cv Canvas) bool {
	wi, hi := cv.Size()
	w, h := float64(wi), float64(hi)
	return p[0] >= 0 && p[0] < w && p[1] >= 0 && p[1] < h
}

// Draw draws the segment on the canvas.
func (s Segment) Draw(cv Canvas, cl color.Color) {
	const length = 25
	s[0].Draw(cv, cl)
	s[1].Draw(cv, cl)
	x0, y0 := int(s[0][0]+0.5), int(s[0][1]+0.5)
	x1, y1 := int(s[1][0]+0.5), int(s[1][1]+0.5)
	cv.StrokeLine(cl, x0, y0, x1, y1)
	s.Normal().ScaledBy(length).DrawAt(cv, cl, s.Center())
}

// Draw draws a circle on the canvas.
func (cir Circle) Draw(cv Canvas, cl color.Color) {
	const N = 100
	const dt = 2 * math.Pi / N

	x0 := int(cir.Center[0] + cir.Radius + 0.5)
	y0 := int(cir.Center[1] + 0.5)
	for i := 1; i < N+1; i++ {
		t := float64(i) * dt
		x1 := int(cir.Center[0] + cir.Radius*math.Cos(t) + 0.5)
		y1 := int(cir.Center[1] + cir.Radius*math.Sin(t) + 0.5)
		cv.StrokeLine(cl, x0, y0, x1, y1)
		x0, y0 = x1, y1
	}
}

// An ImageCanvas implements the Canvas interface using the
// image/draw package from the Go standard library.
type ImageCanvas struct {
	draw.Image
}

// Size returns the size of the canvas in pixels.
func (img ImageCanvas) Size() (int, int) {
	b := img.Bounds()
	return b.Max.X - b.Min.X, b.Max.Y - b.Min.Y
}

// ToImgCoords returns x and y transformed from the Canvas
// frame (0,0 in the lower left, etc.) to the image frame.
func (img ImageCanvas) toImgCoords(x, y int) (int, int) {
	b := img.Bounds()
	h := b.Max.Y - b.Min.Y
	x = x + b.Min.X
	y = h - y - 1 + b.Min.Y
	return x, y
}

// StrokeLine draws a colored line on the canvas.
func (img ImageCanvas) StrokeLine(c color.Color, x0, y0, x1, y1 int) {
	x0, y0 = img.toImgCoords(x0, y0)
	x1, y1 = img.toImgCoords(x1, y1)

	// Bresenham's alg: http://en.wikipedia.org/wiki/Bresenham's_line_algorithm
	steep := abs(y0-y1) > abs(x0-x1)
	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	dx := x1 - x0
	dy := abs(y1 - y0)
	err := dx / 2
	y := y0

	ystep := -1
	if y0 < y1 {
		ystep = 1
	}

	for x := x0; x <= x1; x++ {
		if steep {
			img.Set(y, x, c)
		} else {
			img.Set(x, y, c)
		}
		err -= dy
		if err < 0 {
			y += ystep
			err += dx
		}
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// FillCircle fills a colored circle on the canvas.
func (img ImageCanvas) FillCircle(c color.Color, x, y, r int) {
	p := image.Pt(img.toImgCoords(x, y))
	src := image.NewUniform(c)
	draw.DrawMask(img, img.Bounds(), src, image.ZP, &dot{p, r}, image.ZP, draw.Over)
}

type dot struct {
	p image.Point
	r int
}

func (c dot) ColorModel() color.Model {
	return color.AlphaModel
}

func (c dot) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c dot) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}
