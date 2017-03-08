/*
This file contains functions for drawing texts.
*/

package gfx

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/veandco/go-sdl2/sdl"
)

// DrawString draws a string.
// The y coordinate is the bottom line of the text.
func DrawString(r *sdl.Renderer, s string, x, y int) {
	cr, g, b, a, _ := r.GetDrawColor()

	col := color.NRGBA{cr, g, b, a}
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  newRendererImage(r),
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(s)
}

// renderImage is a draw.Image implementation that targets an sdl.Renderer.
type rendererImage struct {
	// r is the wrapped Renderer
	r *sdl.Renderer

	// bounds to report by Image.Bounds()
	bounds image.Rectangle
}

func newRendererImage(r *sdl.Renderer) draw.Image {
	b := new(sdl.Rect)
	r.GetViewport(b)

	return &rendererImage{
		r:      r,
		bounds: image.Rect(int(b.X), int(b.Y), int(b.X+b.W), int(b.Y+b.H)),
	}
}

func (ri *rendererImage) ColorModel() color.Model {
	return color.NRGBAModel
}

func (ri *rendererImage) Bounds() image.Rectangle {
	return ri.bounds
}

func (ri *rendererImage) At(x, y int) color.Color {
	return color.NRGBA{}
}

func (ri *rendererImage) Set(x, y int, c color.Color) {
	c2 := color.NRGBAModel.Convert(c).(color.NRGBA)
	ri.r.SetDrawColor(c2.R, c2.G, c2.B, c2.A)
	ri.r.DrawPoint(x, y)
}
