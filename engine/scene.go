package engine

import (
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	// physicsPeriod is the period of model recalculation
	physicsPeriod = time.Millisecond * 10 // 100 / sec

	// presentPeriod is the period of the scene presentation
	presentPeriod = time.Millisecond * 20 // 50 FPS
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
func NewScene(r *sdl.Renderer, w, h int) *Scene {
	s := &Scene{
		r:    r,
		quit: make(chan struct{}),
		wg:   &sync.WaitGroup{},
		e:    newEngine(w, h),
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
			sdl.Do(s.present)
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
	r.SetDrawColor(100, 100, 100, 255)
	r.DrawRect(&sdl.Rect{X: 0, Y: 0, W: int32(s.e.w), H: int32(s.e.h)})

	// Paint balls:
	r.SetDrawColor(200, 80, 0, 255)
	for _, b := range s.e.balls {
		paintBall(r, b)
	}

	r.Present()
}

// paintBall paints the picture of a ball, a filled circle with 3D effects.
func paintBall(r *sdl.Renderer, b *ball) {
	// If performance becomes an issue, predraw on a texture,
	// cache it and just present the texture.

	x, y := int(real(b.pos)), int(imag(b.pos))

	// Fill circles going from outside
	gran := 8
	for i := 1; i <= gran; i++ {
		f := 1 - float64(i)/float64(gran+1)

		// Color is darker outside:
		col := func(c uint8) uint8 {
			return c - uint8(float64(c)*0.7*f)
		}

		r.SetDrawColor(col(b.c.R), col(b.c.G), col(b.c.B), b.c.A)
		fillCircle(r, x, y, int(b.r*f))
	}

	r.SetDrawColor(255, 255, 255, b.c.A)
	r.DrawPoint(x, y)
}

// fillCircle draws a filled circle.
func fillCircle(r *sdl.Renderer, x0, y0, rad int) {
	// Algorithm: https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
	for x, y, err := rad, 0, 0; x > 0; {
		r.DrawLine(x0-x, y0-y, x0+x, y0-y)
		r.DrawLine(x0-x, y0+y, x0+x, y0+y)

		if err <= 0 {
			y++
			err += 2*y + 1
		}
		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}
}