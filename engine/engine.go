package engine

import (
	"math/cmplx"
	"time"
)

const (
	// maxBalls is the max number of balls
	maxBalls = 20

	// ballSpawnPeriod is the ball spawning period
	ballSpawnPeriod = time.Second * 2
)

// engine is the simulation engine.
// Contains the balls and simulates the "world".
type engine struct {
	// width and height of the world
	w, h int

	// lastCalc is the last calculation timestamp
	lastCalc time.Time

	// lastSpawned is the last ball spawning timestamp
	lastSpawned time.Time

	// balls of the simulation
	balls []*ball
}

// newEngine creates a new engine.
func newEngine(w, h int) *engine {
	e := &engine{
		w:        w,
		h:        h,
		lastCalc: time.Now(),
	}

	return e
}

// recalc recalculates the world.
func (e *engine) recalc(now time.Time) {
	dt := now.Sub(e.lastCalc)

	if len(e.balls) < maxBalls && now.Sub(e.lastSpawned) > ballSpawnPeriod {
		e.spawnBall()
		e.lastSpawned = now
	}

	dtSec := float64(dt) / float64(time.Second)
	for _, b := range e.balls {
		oldX, oldY := real(b.pos), imag(b.pos)
		b.recalc(dtSec)
		x, y := real(b.pos), imag(b.pos)
		// Check if world boundaries are reached, and bounce back if so:
		if x < b.r || x >= float64(e.w)-b.r {
			b.v = complex(-real(b.v), imag(b.v))
			b.pos = complex(oldX, y)
		}
		if y < b.r || y >= float64(e.h)-b.r {
			b.v = cmplx.Conj(b.v)
			b.pos = complex(x, oldY)
		}
	}

	e.lastCalc = now
}

// spawnBall spawns a new ball.
func (e *engine) spawnBall() {
	// TODO
	b := newBall()

	e.balls = append(e.balls, b)
}
