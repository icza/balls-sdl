package engine

import (
	"math"
	"math/cmplx"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	// gravity is the gravitational constant (vector)
	gravity = 0 + 150i

	// initialMaxAbsV is the max of the absolute of the initial speed vector
	initialMaxAbsV = 300

	// minR is the minimum radius of the ball
	minR = 10

	// maxR is the maximum radius of the ball
	maxR = 30
)

// ball represents a ball.
type ball struct {
	// pos is the position (vector) of the ball
	pos complex128

	// r is the radius of the ball
	r float64

	// v is the velocity (vector) of the ball
	v complex128

	// c is the color of the ball
	c sdl.Color
}

// newBall creates a new ball.
func newBall(w, h int) *ball {
	b := &ball{
		pos: complex(float64(w)/2, float64(h)/2),
		r:   minR + rand.Float64()*(maxR-minR),
		v:   cmplx.Rect(rand.Float64()*initialMaxAbsV, rand.Float64()*math.Pi*2),
		c: sdl.Color{
			R: 127 + uint8(rand.Int31n(128)),
			G: 127 + uint8(rand.Int31n(128)),
			B: 127 + uint8(rand.Int31n(128)),
			A: 255,
		},
	}

	return b
}

// recalc recalculates the position and velocity of the ball based on the delta time.
func (b *ball) recalc(dtSec float64) {
	f := complex(dtSec, 0)

	// simulate gravity:
	b.v += gravity * f

	// Step
	b.pos += b.v * f
}
