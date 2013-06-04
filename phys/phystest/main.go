package main

import (
	"fmt"
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

	// Jerk is the amout by which the acceleration changes when an
	// arrow key is pressed.
	jerk = 5
)

var (
	black  = color.RGBA{A: 255}
	white  = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	red    = color.RGBA{R: 255, A: 255}
	green  = color.RGBA{G: 255, A: 255}
	blue   = color.RGBA{B: 255, A: 255}
	purple = color.RGBA{R: 255, B: 255, A: 255}
	teal   = color.RGBA{G: 255, B: 255, A: 255}
)

var (
	vel, acc Vector
	circle   = Circle{Center: Point{200, 200}, Radius: 50}

	// Sides is the set of polygon sides.
	sides = []Side{
		{{400, 200}, {300, 300}},
		{{300, 300}, {400, 400}},
		{{400, 400}, {500, 300}},
		{{500, 300}, {400, 200}},
	}

	// Click is the position of the latest mouse clicke
	click = Point{-1, -1}

	// Cursor is the current cursor position.
	cursor Point

	// Changed is true if anything has changeds and needs a redraw
	changed bool
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
			if !changed && acc.Equals(Vector{}) {
				continue
			}
			changed = false
			vel.Add(acc)
			circle = phys.MoveCircle(circle, vel, sides)
			vel = Vector{}
			drawScene(win)
		}
	}
}

func mouseMove(ev wde.MouseEvent) {
	cursor = Point{float64(ev.Where.X), float64(height - ev.Where.Y - 1)}
	changed = true
}

func mouseDown(ev wde.MouseButtonEvent) {
	switch ev.Which {
	case wde.LeftButton:
		click = Point{float64(ev.Where.X), float64(height - ev.Where.Y - 1)}
		changed = true
	}
}

func mouseUp(ev wde.MouseButtonEvent) {
	switch ev.Which {
	case wde.LeftButton:
		sides = append(sides, Side{click, cursor})
		click = Point{-1, -1}
		changed = true
	}
}

func keyDown(ev wde.KeyEvent) {
	switch ev.Key {
	case "left_arrow":
		acc[0] = -jerk
	case "right_arrow":
		acc[0] = jerk
	case "up_arrow":
		acc[1] = jerk
	case "down_arrow":
		acc[1] = -jerk
	default:
		fmt.Println("Pressed ", ev.Key)
	}
}

func keyUp(ev wde.KeyEvent) {
	switch ev.Key {
	case "left_arrow", "right_arrow":
		acc[0] = 0
	case "up_arrow", "down_arrow":
		acc[1] = 0
	}
}

func drawScene(win wde.Window) {
	clear(win)
	cv := ImageCanvas{win.Screen()}

	for _, s := range sides {
		s.Draw(cv, black)
	}
	circle.Draw(cv, black)

	if click[0] >= 0 {
		Side{click, cursor}.Draw(cv, blue)
	}

	win.FlushImage()
}

func clear(win wde.Window) {
	img := win.Screen()
	draw.Draw(img, img.Bounds(), image.NewUniform(white), image.ZP, draw.Src)
}
