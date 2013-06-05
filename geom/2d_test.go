package geom

import (
	"image"
	"testing"
)

func TestToImgCoords(t *testing.T) {
	const (
		w = 101
		h = 4000
	)
	tests := []struct {
		x0, y0, x1, y1 int
	}{
		{0, 0, 0, h - 1},
		{w - 1, 0, w - 1, h - 1},
		{w - 1, h - 1, w - 1, 0},
		{0, h - 1, 0, 0},
	}
	c := ImageCanvas{image.NewRGBA(image.Rect(0, 0, w, h))}
	for _, test := range tests {
		x, y := c.toImgCoords(test.x0, test.y0)
		if x == test.x1 && y == test.y1 {
			continue
		}
		t.Errorf("Expect %d,%d to transform to %d,%d, but got %d,%d\n",
			test.x0, test.y0, test.x1, test.y1, x, y)
	}
}

func TestSegmentNormal(t *testing.T) {
	tests := []struct {
		s0, s1 Point
		n      Vector
	}{
		{Point{-1, 0}, Point{1, 0}, Vector{0, 1}},
		{Point{1, 0}, Point{-1, 0}, Vector{0, -1}},
		{Point{0, -1}, Point{0, 1}, Vector{-1, 0}},
		{Point{0, 1}, Point{0, -1}, Vector{1, 0}},
	}

	for _, test := range tests {
		n := Segment{test.s0, test.s1}.Normal()
		if n.Equals(test.n) {
			continue
		}
		t.Errorf("Expect normal of %v to %v to be %v, got %v", test.s0, test.s1, test.n, n)
	}
}
