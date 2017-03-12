package engine

import (
	"math"
	"math/cmplx"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	// gravity is the gravitational constant (vector)
	gravity = 0 + 600i

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

	// m is the mass of the ball
	m complex128 // complex128 for easy use later on
}

// newBall creates a new ball.
func newBall(w, h int) *ball {
	b := &ball{
		pos: complex(
			2*maxR+float64(w-maxR*4)*rand.Float64(),
			float64(h)*0.4,
		),
		//pos: complex(float64(w)*0.5, float64(h)*0.3),
		r: minR + rand.Float64()*(maxR-minR),
		v: cmplx.Rect(rand.Float64()*initialMaxAbsV, rand.Float64()*math.Pi*2),
		c: sdl.Color{
			R: 127 + uint8(rand.Int31n(128)),
			G: 127 + uint8(rand.Int31n(128)),
			B: 127 + uint8(rand.Int31n(128)),
			A: 255,
		},
	}

	// Mass is proportional with the volume which is: V = 4 * r^3 * PI / 3
	b.m = complex(4/3*math.Pi*b.r*b.r*b.r, 0)

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
