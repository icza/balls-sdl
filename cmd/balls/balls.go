// Package main is the Bouncing balls demo app.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/icza/balls/engine"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	version  = "v1.0.0"
	name     = "Bouncing Balls"
	homePage = "https://github.com/icza/balls"
	title    = name + " " + version
)

func main() {
	fmt.Println(title)
	fmt.Println("Home page:", homePage)
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

	// eng is the engine
	eng *engine.Engine
)

// run runs the demo.
func run() (exitCode int) {
	var err error
	fail := func(op string, exitCode int) int {
		log.Printf("could not %s: %v", op, err)
		return exitCode
	}

	sdl.Do(func() {
		err = sdl.Init(sdl.INIT_VIDEO)
	})
	if err != nil {
		return fail("init SDL video", 1)
	}
	defer sdl.Do(sdl.Quit)

	bounds := new(sdl.Rect)
	sdl.Do(func() {
		err = sdl.GetDisplayBounds(0, bounds)
	})
	if err != nil {
		return fail("get display bounds", 2)
	}
	// Start with a half-size window
	w, h := int(bounds.W)/2, int(bounds.H)/2

	sdl.Do(func() {
		win, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, w, h, sdl.WINDOW_SHOWN)
	})
	if err != nil {
		return fail("create window", 3)
	}
	defer sdl.Do(win.Destroy)

	var r *sdl.Renderer
	sdl.Do(func() {
		r, err = sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED)
	})
	if err != nil {
		return fail("create renderer", 4)
	}
	defer sdl.Do(r.Destroy)

	sdl.Do(func() {
		// set logical size, so if window gets resized (full screen),
		// the world size does not change:
		err = r.SetLogicalSize(w, h)
	})
	if err != nil {
		return fail("set logical size", 5)
	}

	sdl.Do(func() {
		// Disable minimize on focus loss:
		if !sdl.SetHint(sdl.HINT_VIDEO_MINIMIZE_ON_FOCUS_LOSS, "0") {
			log.Println("Waring: Failed to disable HINT_VIDEO_MINIMIZE_ON_FOCUS_LOSS!")
		}
	})

	eng = engine.NewEngine(r, w, h)
	go eng.Run()

	for {
		var e sdl.Event
		sdl.Do(func() {
			e = sdl.PollEvent()
		})
		if e != nil {
			if quit := handleEvent(e); quit {
				break
			}
		} else {
			time.Sleep(time.Millisecond)
		}
	}
	eng.Stop()

	return 0
}

// handleEvent handles events and tells if we need to quit (based on the event).
func handleEvent(event sdl.Event) (quit bool) {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_f:
			if time.Since(lastFSSwitch) > time.Millisecond*500 {
				flags := uint32(0)
				if !fullScreen {
					flags = sdl.WINDOW_FULLSCREEN_DESKTOP
				}
				sdl.Do(func() {
					win.SetFullscreen(flags)
				})
				fullScreen = !fullScreen
				lastFSSwitch = time.Now()
			}
		case sdl.K_s:
			eng.ChangeSpeed(e.Keysym.Mod&sdl.KMOD_SHIFT != 0)
		case sdl.K_a:
			eng.ChangeMaxBalls(e.Keysym.Mod&sdl.KMOD_SHIFT != 0)
		case sdl.K_m:
			eng.ChangeMinMaxBallRatio(e.Keysym.Mod&sdl.KMOD_SHIFT != 0)
		case sdl.K_r:
			eng.Restart()
		case sdl.K_o:
			eng.ToggleOSD()
		case sdl.K_g:
			eng.ChangeGravityAbs(e.Keysym.Mod&sdl.KMOD_SHIFT != 0)
		case sdl.K_t:
			eng.RotateGravity(e.Keysym.Mod&sdl.KMOD_SHIFT != 0)
		case sdl.K_x, sdl.K_q:
			return true
		}
	case *sdl.QuitEvent:
		return true
	}

	return false
}
