package geom

import (
	"image"
	"testing"
)

func TestImageCanvas_toImgCoords(t *testing.T) {
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
