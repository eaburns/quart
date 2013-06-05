package main

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"time"

	. "github.com/eaburns/quart/geom"
	"github.com/eaburns/quart/phys"

	"github.com/skelterjohn/go.wde"
)

const (
	width  = 640
	height = 480

	speed         = 5
	gravity       = -3
	stopThreshold = 1
)

var (
	vel    Vector
	circle = Circle{Center: Point{200, 200}, Radius: 50}

	// Segs is the set of segments defining obstacles.
	segs = []Segment{
		{{0, height - 1}, {0, 0}},
		{{0, 0}, {width - 1, 0}},
		{{width - 1, 0}, {width - 1, height - 1}},
		{{width - 1, height - 1}, {0, height - 1}},
	}

	// Click is the position of the latest mouse click.
	click = Point{-1, -1}

	// Cursor is the current cursor position.
	cursor Point

	// Stopped is true if the circle has effectively stopped moving.
	stopped bool
)

func main() {
	go mainLoop()
	wde.Run()
}

func mainLoop() {
	win, err := wde.NewWindow(width, height)
	if err != nil {
		panic(err)
	}
	win.SetTitle("geom test")
	win.Show()

	drawScene(win)

	tick := time.NewTicker(40 * time.Millisecond)
	for {
		select {
		case ev, ok := <-win.EventChan():
			if !ok {
				os.Exit(0)
			}
			switch ev := ev.(type) {
			case wde.CloseEvent:
				os.Exit(0)
			case wde.KeyTypedEvent:
				keyTyped(ev)
			case wde.KeyDownEvent:
				keyDown(wde.KeyEvent(ev))
			case wde.KeyUpEvent:
				keyUp(wde.KeyEvent(ev))
			case wde.MouseDraggedEvent:
				mouseMove(ev.MouseEvent)
			case wde.MouseMovedEvent:
				mouseMove(ev.MouseEvent)
			case wde.MouseDownEvent:
				mouseDown(wde.MouseButtonEvent(ev))
			case wde.MouseUpEvent:
				mouseUp(wde.MouseButtonEvent(ev))
			}

		case <-tick.C:
			if !stopped {
				start := circle.Center
				circle = phys.MoveCircle(circle, vel.Plus(Vector{0, gravity}), segs)
				dist := start.Minus(circle.Center).Magnitude()
				stopped = vel.Equals(Vector{}) && dist < stopThreshold
			}
			drawScene(win)
		}
	}
}

func mouseMove(ev wde.MouseEvent) {
	cursor = Point{float64(ev.Where.X), float64(height - ev.Where.Y - 1)}
}

func mouseDown(ev wde.MouseButtonEvent) {
	switch ev.Which {
	case wde.LeftButton:
		click = Point{float64(ev.Where.X), float64(height - ev.Where.Y - 1)}
	}
}

func mouseUp(ev wde.MouseButtonEvent) {
	switch ev.Which {
	case wde.LeftButton:
		segs = append(segs, Segment{click, cursor})
		click = Point{-1, -1}
	}
}

func keyTyped(ev wde.KeyTypedEvent) {
	switch ev.Key {
	case "d":
		if len(segs) > 4 {
			segs = segs[:len(segs)-1]
		}
	}
}

func keyDown(ev wde.KeyEvent) {
	switch ev.Key {
	case "left_arrow":
		vel[0] = -speed
	case "right_arrow":
		vel[0] = speed
	case "up_arrow":
		vel[1] = speed
	case "down_arrow":
		vel[1] = -speed
	}
	stopped = false
}

func keyUp(ev wde.KeyEvent) {
	switch ev.Key {
	case "left_arrow", "right_arrow":
		vel[0] = 0
	case "up_arrow", "down_arrow":
		vel[1] = 0
	}
}

func drawScene(win wde.Window) {
	clear(win)
	cv := ImageCanvas{win.Screen()}

	for _, s := range segs {
		s.Draw(cv, color.Black)
	}
	circle.Draw(cv, color.Black)

	if click[0] >= 0 {
		Segment{click, cursor}.Draw(cv, color.RGBA{B: 255, A: 255})
	}

	win.FlushImage()
}

func clear(win wde.Window) {
	img := win.Screen()
	draw.Draw(img, img.Bounds(), image.NewUniform(color.White), image.ZP, draw.Src)
}
