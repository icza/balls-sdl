// Package main is the Bouncing balls demo app.
package main

import (
	"log"
	"os"

	"github.com/icza/balls/engine"
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
	fail := func(op string, exitCode int) int {
		log.Printf("could not %s: %v", op, err)
		return exitCode
	}

	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		return fail("init SDL video", 1)
	}
	defer sdl.Quit()

	var win *sdl.Window
	if win, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 800, 600, sdl.WINDOW_SHOWN); err != nil {
		return fail("create window", 2)
	}
	defer win.Destroy()

	var r *sdl.Renderer
	if r, err = sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED); err != nil {
		return fail("create renderer", 3)
	}
	defer r.Destroy()

	scene := engine.NewScene(r)
	go scene.Run()

	for {
		e := sdl.WaitEvent()
		if quit := handleEvent(e); quit {
			break
		}
	}
	scene.Stop()

	return 0
}

func handleEvent(event sdl.Event) (quit bool) {
	switch e := event.(type) {
	case *sdl.QuitEvent:
		return true
	// Ignored events:
	case *sdl.MouseMotionEvent:
	default:
		log.Printf("event: %T", e)
	}

	return false
}
