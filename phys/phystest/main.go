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
	acc   = Vector{}
	vel   = Vector{}
	sides = []Side{
		{{300, 300}, {400, 200}},
		{{400, 400}, {300, 300}},
		{{500, 300}, {400, 400}},
		{{400, 200}, {500, 300}},
	}
	circle = Circle{Center: Point{200, 200}, Radius: 50}
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

			case wde.KeyUpEvent:
				switch ev.Key {
				case "left_arrow", "right_arrow":
					acc[0] = 0
				case "up_arrow", "down_arrow":
					acc[1] = 0
				}
			}

		case <-tick.C:
			if acc.Equals(Vector{}) {
				continue
			}
			vel.Add(acc)
			circle = phys.MoveCircle(circle, vel, sides)
			vel = Vector{}
			drawScene(win)
		}
	}
}

func drawScene(win wde.Window) {
	clear(win)
	cv := ImageCanvas{win.Screen()}

	for _, s := range sides {
		s.Draw(cv, blue)
	}
	circle.Draw(cv, black)
	win.FlushImage()
}

func clear(win wde.Window) {
	img := win.Screen()
	draw.Draw(img, img.Bounds(), image.NewUniform(white), image.ZP, draw.Src)
}
