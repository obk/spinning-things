package main

import (
	"bufio"
	"math"
)

// The cube is drawn with one character and color per face, using a z-buffer
// so nearer surfaces hide the ones behind them.

const (
	cubeHalf     = 10.0 // half the cube's edge length
	cubeStep     = 0.6  // sampling step across each face
	cubeDistance = 60.0 // distance from camera to cube center
	cubeScale    = 40.0 // projection scale
)

type face struct {
	char  byte
	color string
}

// One entry per face: front, right, back, left, bottom, top.
var faces = [6]face{
	{'@', "\x1b[31m"},
	{'$', "\x1b[32m"},
	{'~', "\x1b[33m"},
	{'#', "\x1b[34m"},
	{';', "\x1b[35m"},
	{'+', "\x1b[36m"},
}

type cube struct {
	a, b, c float64 // rotation angles around x, y, z
	zbuf    [width * height]float64
	buf     [width * height]int8 // face index per cell, -1 = empty
}

// rotate applies rotations around the x, y and z axes to a point.
func (r *cube) rotate(x, y, z float64) (float64, float64, float64) {
	sa, ca := math.Sincos(r.a)
	sb, cb := math.Sincos(r.b)
	sc, cc := math.Sincos(r.c)

	// Rotate around X.
	y, z = y*ca-z*sa, y*sa+z*ca
	// Rotate around Y.
	x, z = x*cb+z*sb, -x*sb+z*cb
	// Rotate around Z.
	x, y = x*cc-y*sc, x*sc+y*cc
	return x, y, z
}

// plot projects a 3D point onto the screen and records it if it is nearest.
func (r *cube) plot(x, y, z float64, f int8) {
	x, y, z = r.rotate(x, y, z)
	z += cubeDistance
	ooz := 1 / z // "one over z": larger means closer

	// Terminal cells are roughly twice as tall as they are wide, so x is doubled.
	xp := int(width/2 + cubeScale*ooz*x*2)
	yp := int(height/2 + cubeScale*ooz*y)
	if xp < 0 || xp >= width || yp < 0 || yp >= height {
		return
	}
	idx := yp*width + xp
	if ooz > r.zbuf[idx] {
		r.zbuf[idx] = ooz
		r.buf[idx] = f
	}
}

func (r *cube) frame(out *bufio.Writer) {
	for i := range r.buf {
		r.buf[i] = -1
		r.zbuf[i] = 0
	}

	const h = cubeHalf
	for u := -h; u < h; u += cubeStep {
		for v := -h; v < h; v += cubeStep {
			r.plot(u, v, -h, 0)  // front
			r.plot(h, v, u, 1)   // right
			r.plot(-u, v, h, 2)  // back
			r.plot(-h, v, -u, 3) // left
			r.plot(u, -h, -v, 4) // bottom
			r.plot(u, h, v, 5)   // top
		}
	}

	last := int8(-1)
	for y := range height {
		for x := range width {
			f := r.buf[y*width+x]
			if f < 0 {
				out.WriteByte(' ')
				continue
			}
			if f != last {
				out.WriteString(faces[f].color)
				last = f
			}
			out.WriteByte(faces[f].char)
		}
		out.WriteByte('\n')
	}

	r.a += 0.05
	r.b += 0.05
	r.c += 0.01
}
