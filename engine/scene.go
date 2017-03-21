package engine

import (
	"fmt"

	"github.com/icza/balls/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

// scene is used to present the world.
type scene struct {
	// r is the Renderer used to paint.
	r *sdl.Renderer

	// e is the engine
	e *Engine
}

// newScene creates a new scene.
func newScene(r *sdl.Renderer, e *Engine) *scene {
	s := &scene{
		r: r,
		e: e,
	}

	return s
}

// present paints the scene in the SDL2's main thread.
func (s *scene) present() {
	sdl.Do(s.presentInternal)
}

// presentInternal paints the scene.
func (s *scene) presentInternal() {
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

	s.paintOSD()

	r.Present()
}

// paintOSD paints on-screen texts.
func (s *scene) paintOSD() {
	if !s.e.osd {
		return
	}

	r := s.r

	r.SetDrawColor(200, 200, 100, 255)
	speed := 1.0
	if exp := s.e.speedExp; exp >= 0 {
		speed *= float64(int(1) << uint(exp))
	} else {
		speed /= float64(int(1) << uint(-exp))
	}

	items := []struct {
		keys   string
		format string
		param  interface{}
	}{
		{"F", "fullscreen", nil},
		{"R", "restart", nil},
		{"O", "Toggle OSD (help text)", nil},
		{"Q/X", "Quit", nil},
		{"S/s", "speed: %.2f", speed},
		{"A/a", "max # of balls: %2d", s.e.maxBalls},
		{"M/m", "min-max ball ratio: %.1f", float64(s.e.minMaxBallRatio) / 100},
	}

	row, col := 0, 0
	xf := func(col int) int { return col*210 + 10 }
	for _, it := range items {
		params := []interface{}{"[" + it.keys + "]"}
		if it.param != nil {
			params = append(params, it.param)
		}
		text := fmt.Sprintf("%-5s "+it.format, params...)
		gfx.DrawString(r, text, xf(col), row*15+15)

		col++
		if xf(col+1) > s.e.w {
			row, col = row+1, 0
		}
	}
}

// paintBall paints the picture of a ball, a filled circle with 3D effects.
func paintBall(r *sdl.Renderer, b *ball) {
	// If performance becomes an issue, predraw on a texture,
	// cache it and just present the texture.

	x, y := int(real(b.pos)), int(imag(b.pos))

	// Fill circles going from outside
	gran := 10
	for i := 1; i <= gran; i++ {
		f := 1 - float64(i)/float64(gran+1)

		// Color is darker outside:
		col := func(c uint8) uint8 {
			return c - uint8(float64(c)*0.7*f)
		}

		r.SetDrawColor(col(b.c.R), col(b.c.G), col(b.c.B), b.c.A)
		gfx.FillCircle(r, x, y, int(b.r*f))
	}

	r.SetDrawColor(255, 255, 255, b.c.A)
	r.DrawPoint(x, y)
}
