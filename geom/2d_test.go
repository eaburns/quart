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
		t.Errorf("Expected %d,%d to transform to %d,%d, but got %d,%d\n",
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
		if n.NearlyEquals(test.n) {
			continue
		}
		t.Errorf("Expected normal of %v to %v to be %v, got %v", test.s0, test.s1, test.n, n)
	}
}

func TestRectangleMax(t *testing.T) {
	tests := []struct {
		min  Point
		size Vector
		max  Point
	}{
		{Point{0, 0}, Vector{1, 1}, Point{1, 1}},
		{Point{0, 1}, Vector{1, 1}, Point{1, 2}},
	}
	for _, test := range tests {
		r := Rectangle{test.min, test.size}
		m := r.Max()
		if m.NearlyEquals(test.max) {
			continue
		}
		t.Errorf("Expected max of %v to be %v, got %v", r, test.max, m)
	}
}

func TestRectangleCenter(t *testing.T) {
	tests := []struct {
		min    Point
		size   Vector
		center Point
	}{
		{Point{0, 0}, Vector{1, 1}, Point{0.5, 0.5}},
		{Point{-1, -1}, Vector{2, 2}, Point{0, 0}},
	}
	for _, test := range tests {
		r := Rectangle{test.min, test.size}
		c := r.Center()
		if c.NearlyEquals(test.center) {
			continue
		}
		t.Errorf("Expected center of %v to be %v, got %v", r, test.center, c)
	}
}

func BenchmarkLineDirection(b *testing.B) {
	l := Line{Origin: Point{0, 0}, Normal: Vector{0, 1}}
	for i := 0; i < b.N; i++ {
		l.Direction()
	}
}

func BenchmarkLineLineIntersectionHit(b *testing.B) {
	l0 := Line{Origin: Point{0, 0}, Normal: Vector{0, 1}}
	l1 := Line{Origin: Point{0, 0}, Normal: Vector{1, 0}}
	for i := 0; i < b.N; i++ {
		l0.LineIntersection(l1)
	}
}

func BenchmarkLineLineIntersectionMiss(b *testing.B) {
	l0 := Line{Origin: Point{0, 0}, Normal: Vector{0, 1}}
	l1 := Line{Origin: Point{1, 0}, Normal: Vector{0, 1}}
	for i := 0; i < b.N; i++ {
		l0.LineIntersection(l1)
	}
}

func BenchmarkSegmentNormal(b *testing.B) {
	s := Segment{Point{0, 0}, Point{1, 0}}
	for i := 0; i < b.N; i++ {
		s.Normal()
	}
}

func BenchmarkSegmentLine(b *testing.B) {
	s := Segment{Point{0, 0}, Point{1, 0}}
	for i := 0; i < b.N; i++ {
		s.Line()
	}
}

func BenchmarkRectangleMax(b *testing.B) {
	r := Rectangle{Point{0, 0}, Vector{1, 1}}
	for i := 0; i < b.N; i++ {
		r.Max()
	}
}

func BenchmarkRectangleCenter(b *testing.B) {
	r := Rectangle{Point{0, 0}, Vector{1, 1}}
	for i := 0; i < b.N; i++ {
		r.Center()
	}
}
