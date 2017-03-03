package engine

import "time"

const (
	// maxBalls is the max number of balls
	maxBalls = 20

	// ballSpawnPeriod is the ball spawning period
	ballSpawnPeriod = time.Second * 2
)

// engine is the simulation engine.
// Contains the balls and simulates the "world".
type engine struct {
	// lastCalc is the last calculation timestamp
	lastCalc time.Time

	// lastSpawned is the last ball spawning timestamp
	lastSpawned time.Time

	// balls of the simulation
	balls []*ball
}

// newEngine creates a new engine.
func newEngine() *engine {
	e := &engine{
		lastCalc: time.Now(),
	}

	return e
}

// recalc recalculates the world.
func (e *engine) recalc(now time.Time) {
	dt := now.Sub(e.lastCalc)

	if now.Sub(e.lastSpawned) > ballSpawnPeriod {
		e.spawnBall()
		e.lastSpawned = now
	}

	dtSec := float64(dt) / float64(time.Second)
	for _, b := range e.balls {
		b.recalc(dtSec)
	}

	e.lastCalc = now
}

// spawnBall spawns a new ball.
func (e *engine) spawnBall() {
	// TODO
	b := newBall()

	e.balls = append(e.balls, b)
}
