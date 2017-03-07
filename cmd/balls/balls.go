// Package main is the Bouncing balls demo app.
package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/icza/balls/engine"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	title = "Bouncing Balls"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	os.Exit(run())
}

var (
	// win is the main window
	win *sdl.Window
	// fullScreen tells the current FS status
	fullScreen bool
	// lastFSSwitch holds the last FS switch timestamp (to limit the switching rate)
	lastFSSwitch time.Time
)

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

	bounds := new(sdl.Rect)
	if err = sdl.GetDisplayBounds(0, bounds); err != nil {
		return fail("get display bounds", 2)
	}
	// Start with a half-size window
	w, h := int(bounds.W)/2, int(bounds.H)/2

	if win, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, w, h, sdl.WINDOW_SHOWN); err != nil {
		return fail("create window", 3)
	}
	defer win.Destroy()

	var r *sdl.Renderer
	if r, err = sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED); err != nil {
		return fail("create renderer", 4)
	}
	defer r.Destroy()

	// set logical size, so if window gets resized (full screen),
	// the world size does not chane:
	r.SetLogicalSize(w, h)

	scene := engine.NewScene(r, w, h)
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
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_f:
			if time.Since(lastFSSwitch) > time.Second {
				flags := uint32(0)
				if !fullScreen {
					flags = sdl.WINDOW_FULLSCREEN_DESKTOP
				}
				win.SetFullscreen(flags)
				fullScreen = !fullScreen
				lastFSSwitch = time.Now()
			}
		case sdl.K_x:
			return true
		}
	case *sdl.QuitEvent:
		return true
	// Ignored events:
	case *sdl.MouseMotionEvent:
	default:
		log.Printf("event: %T", e)
	}

	return false
}
