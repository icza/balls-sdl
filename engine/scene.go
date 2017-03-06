package engine

import (
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	gfx "github.com/veandco/go-sdl2/sdl_gfx"
)

const (
	// physicsPeriod is the period of model recalculation
	physicsPeriod = time.Millisecond * 10 // 100 / sec

	// presentPeriod is the period of the scene presentation
	presentPeriod = time.Millisecond * 40 // 25 FPS
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
	w, h, _ := r.GetRendererOutputSize()

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
	w, h, _ := r.GetRendererOutputSize()
	r.SetDrawColor(150, 150, 150, 255)
	r.DrawRect(&sdl.Rect{X: 0, Y: 0, W: int32(w), H: int32(h)})

	// Paint balls:
	r.SetDrawColor(200, 80, 0, 255)
	for _, b := range s.e.balls {
		paintBall(r, b)
	}

	r.Present()
}

func paintBall(r *sdl.Renderer, b *ball) {
	// If performance becomes an issue, predraw on a texture,
	// cache it and just present the texture.

	x, y := int(real(b.pos)), int(imag(b.pos))

	// Fill circles going from outside
	gran := 7
	for i := 1; i <= gran; i++ {
		f := 1 - float64(i)/float64(gran+1)

		// Color is darker outside:
		col := func(c uint8) uint8 {
			return c - uint8(float64(c)*0.7*f)
		}

		gfx.FilledCircleRGBA(r, x, y, int(b.r*f),
			col(b.c.R), col(b.c.G), col(b.c.B), b.c.A)
	}

	gfx.PixelRGBA(r, x, y, 255, 255, 255, b.c.A)
}
