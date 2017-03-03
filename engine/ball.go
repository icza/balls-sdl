package engine

import (
	"math"
	"math/cmplx"
	"math/rand"
)

const (
	// gravity is the gravitational constant (vector)
	gravity = 0 + 120i

	// initialMaxAbsV is the max of the absolute of the initial speed vector
	initialMaxAbsV = 100

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
}

// newBall creates a new ball.
func newBall() *ball {
	// TODO
	v := cmplx.Rect(rand.Float64()*initialMaxAbsV, rand.Float64()*math.Pi*2)

	b := &ball{
		pos: 400 + 300i,
		r:   minR + rand.Float64()*(maxR-minR),
		v:   v,
	}

	return b
}

// recalc recalculates the position and velocity of the ball based on the delta time.
func (b *ball) recalc(dtSec float64) {
	f := complex(dtSec, 0)

	// Step
	b.pos += b.v * f

	// simulate gravity:
	b.v += gravity * f
}
