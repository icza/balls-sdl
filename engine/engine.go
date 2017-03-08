package engine

import (
	"math/cmplx"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	// physicsPeriod is the period of model recalculation
	physicsPeriod = time.Millisecond * 10 // 100 / sec

	// presentPeriod is the period of the scene presentation
	presentPeriod = time.Millisecond * 20 // 50 FPS

	// maxBalls is the max number of balls
	maxBalls = 20

	// ballSpawnPeriod is the ball spawning period
	ballSpawnPeriod = time.Second * 2
)

// Engine is the simulation engine.
// Contains the model, controls the simulation and presents it on the screen
// (via the scene).
type Engine struct {
	// w and h are the width and height of the world
	w, h int

	// scene is used to present the world
	scene *scene

	// quit is used to signal termination
	quit chan struct{}

	// wg is a WaitGroup to wait Run to return
	wg *sync.WaitGroup

	// lastCalc is the last calculation timestamp
	lastCalc time.Time

	// lastSpawned is the last ball spawning timestamp
	lastSpawned time.Time

	// balls of the simulation
	balls []*ball
}

// NewEngine creates a new Engine.
func NewEngine(r *sdl.Renderer, w, h int) *Engine {
	e := &Engine{
		w:        w,
		h:        h,
		quit:     make(chan struct{}),
		wg:       &sync.WaitGroup{},
		lastCalc: time.Now(),
	}
	e.scene = newScene(r, e)

	// Add one here (and not in Run()) because if Stop() is called before
	// Run() could start, Stop() would return immediately even though Run()
	// might be started after that.
	e.wg.Add(1)

	return e
}

// Run runs the simulation.
func (e *Engine) Run() {
	defer e.wg.Done()

	ticker := time.NewTicker(presentPeriod)
	defer ticker.Stop()

simLoop:
	for {
		select {
		case now := <-ticker.C:
			e.recalc(now)
			e.scene.present()
		case <-e.quit:
			break simLoop
		}
	}
}

// Stop stops the simulation and waits for Run to return.
func (e *Engine) Stop() {
	close(e.quit)
	e.wg.Wait()
}

// recalc recalculates the world.
func (e *Engine) recalc(now time.Time) {
	// dt might be "big", much bigger than physics period, in which case
	// the balls might move huge distances. To avoid any "unexpected" states,
	// do granular movement:
	for t := e.lastCalc; t.Before(now); t = t.Add(physicsPeriod) {
		e.recalcInternal(t)
	}
}

// recalcInternal recalculates the world.
func (e *Engine) recalcInternal(now time.Time) {
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
func (e *Engine) spawnBall() {
	b := newBall(e.w, e.h)

	// TODO check if no collision

	e.balls = append(e.balls, b)
}
