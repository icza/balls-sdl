package engine

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/icza/balls-sdl/gfx"
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
		s.paintBall(b)
	}

	s.paintOSD()

	s.paintGravity()

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

	phase := cmplx.Phase(s.e.gravity) + math.Pi*2 // Make sure it's positive
	degree := (720 - int(phase/math.Pi*180+0.5)) % 360

	items := []struct {
		keys   string
		format string
		param  interface{}
	}{
		{"F", "fullscreen", nil},
		{"R", "restart", nil},
		{"Q/X", "quit", nil},
		{"O", "OSD (on-screen display)", nil},
		{"S/s", "speed: %.2f", speed},
		{"A/a", "max # of balls: %2d", s.e.maxBalls},
		{"G/g", "abs gravity: %.2f", cmplx.Abs(s.e.gravity) / maxAbsGravity},
		{"T/t", "rotate gravity: %3d deg", degree},
		{"M/m", "min/max ball ratio: %.1f", float64(s.e.minMaxBallRatio) / 100},
	}

	col2x := func(col int) int { return col*210 + 10 }
	row2y := func(row int) int { return row*15 + 15 }

	// How many text columns fits on the screen?
	numCol := 0
	for col2x(numCol+1) < int(s.e.w) {
		numCol++
	}

	row, col := 0, 0
	for _, it := range items {
		params := []interface{}{"[" + it.keys + "]"}
		if it.param != nil {
			params = append(params, it.param)
		}
		text := fmt.Sprintf("%-5s "+it.format, params...)
		gfx.DrawString(r, text, col2x(col), row2y(row))

		col++
		if col >= numCol {
			row, col = row+1, 0
		}
	}
}

// paintBall paints the picture of a ball, a filled circle with 3D effects.
func (s *scene) paintBall(b *ball) {
	// If performance becomes an issue, predraw on a texture,
	// cache it and just present the texture.
	r := s.r

	x, y := int32(real(b.pos)), int32(imag(b.pos))

	// Fill circles going from outside
	gran := 10
	for i := 1; i <= gran; i++ {
		f := 1 - float64(i)/float64(gran+1)

		// Color is darker outside:
		col := func(c uint8) uint8 {
			return c - uint8(float64(c)*0.7*f)
		}

		r.SetDrawColor(col(b.c.R), col(b.c.G), col(b.c.B), b.c.A)
		gfx.FillCircle(r, x, y, int32(b.r*f))
	}

	r.SetDrawColor(255, 255, 255, b.c.A)
	r.DrawPoint(x, y)
}

// paintGravity paints a gravity vector.
func (s *scene) paintGravity() {
	const size = 70 // Pixel size of max gravity
	g := s.e.gravity * complex(float64(size)/maxAbsGravity, 0)

	x1, y1 := int32(s.e.w-size-2), int32(s.e.h-size-2)
	x2, y2 := x1+int32(real(g)), y1+int32(imag(g))

	s.r.SetDrawColor(50, 150, 255, 255)
	s.r.DrawLine(x1, y1, x2, y2)

	// Bottom of the arrow:
	v := g * 0.15i
	s.r.DrawLine(x1, y1, x1+int32(real(v)), y1+int32(imag(v)))
	v = g * -0.15i
	s.r.DrawLine(x1, y1, x1+int32(real(v)), y1+int32(imag(v)))

	// Head of the arrow:
	v = g * (-0.18 + 0.18i)
	s.r.DrawLine(x2, y2, x2+int32(real(v)), y2+int32(imag(v)))
	v = g * (-0.18 - 0.18i)
	s.r.DrawLine(x2, y2, x2+int32(real(v)), y2+int32(imag(v)))
}
