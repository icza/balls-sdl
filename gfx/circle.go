/*
This file contains functions for drawing and filling a circle.

The Midpoint circle algorithm is used which is detailed here:
https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
*/

package gfx

import "github.com/veandco/go-sdl2/sdl"

// FillCircle draws a filled circle.
func FillCircle(r *sdl.Renderer, x0, y0, rad int32) {
	for x, y, err := rad, int32(0), int32(0); x > 0; {
		r.DrawLine(x0-x, y0-y, x0+x, y0-y)
		r.DrawLine(x0-x, y0+y, x0+x, y0+y)

		if err <= 0 {
			y++
			err += 2*y + 1
		}
		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}
}
