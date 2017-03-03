package engine

import (
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	// physicsPeriod is the period of model recalculation
	physicsPeriod = time.Millisecond * 10

	// presentPeriod is the period of the scene presentation
	presentPeriod = time.Millisecond * 50 // 20 FPS
)

// Scene is the world of the demo.
// Contains the engine, controls the simulation and presents it on the screen.
type Scene struct {
	// r is the Renderer used to paint.
	r *sdl.Renderer

	// quit is used to signal termination
	quit chan struct{}

	// wg is a WaitGroup to wait Run to return
	wg *sync.WaitGroup

	// e is the engine
	e *engine
}

// NewScene creates a new Scene.
func NewScene(r *sdl.Renderer) *Scene {
	s := &Scene{
		r:    r,
		quit: make(chan struct{}),
		wg:   &sync.WaitGroup{},
		e:    newEngine(),
	}

	// Add one here (and not in Run()) because if Stop() is called before
	// Run() could start, Stop() would return immediately even though Run()
	// might be started after that.
	s.wg.Add(1)

	return s
}

// Run runs the simulation.
func (s *Scene) Run() {
	defer s.wg.Done()

	physicsTicker := time.NewTicker(physicsPeriod)
	defer physicsTicker.Stop()
	presentTicker := time.NewTicker(presentPeriod)
	defer presentTicker.Stop()

simLoop:
	for {
		select {
		case now := <-physicsTicker.C:
			s.e.recalc(now)
		case <-presentTicker.C:
			s.present()
		case <-s.quit:
			break simLoop
		}
	}
}

// Stop stops the simulation and waits for Run to return.
func (s *Scene) Stop() {
	close(s.quit)
	s.wg.Wait()
}

// present paints the scene.
func (s *Scene) present() {
	r := s.r

	r.SetDrawColor(0, 0, 0, 255)
	r.Clear()

	// Paint background and frame:
	r.SetDrawColor(150, 150, 150, 255)
	r.DrawRect(&sdl.Rect{X: 0, Y: 0, W: 800, H: 600})

	// Paint balls:
	r.SetDrawColor(200, 80, 0, 255)
	for _, b := range s.e.balls {
		_ = b
		//gfx.FilledCircleRGBA()
		r.DrawRect(&sdl.Rect{
			X: int32(real(b.pos) - b.r),
			Y: int32(imag(b.pos) - b.r),
			W: int32(b.r * 2),
			H: int32(b.r * 2),
		})
	}

	r.Present()
}
