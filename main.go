// Package main is the runnable app of the Bouncing balls demo app.
package main

import (
	"log"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Printf("could not init SDL video: %v", err)
		return 1
	}
	defer sdl.Quit()

	return 0
}
