// Package main is the runnable app of the Bouncing balls demo app.
package main

import (
	"log"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	title = "Bouncing Balls"
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	var err error
	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		return fail("init SDL video", err, 1)
	}
	defer sdl.Quit()

	var win *sdl.Window
	if win, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 800, 600, sdl.WINDOW_SHOWN); err != nil {
		return fail("create window", err, 2)
	}
	defer win.Destroy()

	var r *sdl.Renderer
	if r, err = sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED); err != nil {
		return fail("create renderer", err, 3)
	}
	defer r.Destroy()

	r.Clear()
	r.SetDrawColor(150, 150, 150, 255)
	r.DrawRect(&sdl.Rect{X: 0, Y: 0, W: 800, H: 600})
	r.Present()
	sdl.Delay(3000)

	return 0
}

func fail(op string, err error, exitCode int) int {
	log.Printf("could not %s: %v", op, err)
	return exitCode
}
