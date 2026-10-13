package main

import (
	"bufio"
	"math"
)

// The lemniscate of Gerono is the figure-eight curve (cos t, sin t cos t).
// It is drawn as a tube swept along that curve and lit like the donut.

const (
	lemR         = 2.0  // half the figure-eight's width
	lemTube      = 0.35 // radius of the tube
	lemDistance  = 7.0  // distance from camera to curve center
	lemTStep     = 0.01 // step along the curve
	lemThetaStep = 0.05 // step around the tube's cross-section
)

// Projection scale chosen so the figure fills about 3/4 of the screen height.
const lemScale = height * lemDistance * 3 / (8 * (lemR + lemTube))

type lemniscate struct {
	a, b float64 // rotation angles around x and y
	zbuf [width * height]float64
	buf  [width * height]byte
}

// rotate applies rotations around the x and y axes to a vector.
func (l *lemniscate) rotate(x, y, z, sa, ca, sb, cb float64) (float64, float64, float64) {
	y, z = y*ca-z*sa, y*sa+z*ca
	x, z = x*cb+z*sb, -x*sb+z*cb
	return x, y, z
}

func (l *lemniscate) frame(out *bufio.Writer) {
	for i := range l.buf {
		l.buf[i] = ' '
		l.zbuf[i] = 0
	}

	sa, ca := math.Sincos(l.a)
	sb, cb := math.Sincos(l.b)

	// Light from above and behind the viewer, normalized.
	const lx, ly, lz = 0, math.Sqrt2 / 2, -math.Sqrt2 / 2

	for t := 0.0; t < 2*math.Pi; t += lemTStep {
		sinT, cosT := math.Sincos(t)
		// Point on the curve and its unit tangent.
		cx, cy := lemR*cosT, lemR*sinT*cosT
		tx, ty := -sinT, math.Cos(2*t)
		tl := math.Hypot(tx, ty)
		tx, ty = tx/tl, ty/tl
		// The curve is flat, so the binormal is z and the normal lies in the plane.
		nx, ny := -ty, tx

		for theta := 0.0; theta < 2*math.Pi; theta += lemThetaStep {
			sinH, cosH := math.Sincos(theta)
			// Surface normal of the tube at this point.
			ox, oy, oz := cosH*nx, cosH*ny, sinH

			x, y, z := l.rotate(cx+lemTube*ox, cy+lemTube*oy, lemTube*oz, sa, ca, sb, cb)
			z += lemDistance
			ooz := 1 / z // "one over z": larger means closer

			// Terminal cells are roughly twice as tall as they are wide, so x is doubled.
			xp := int(width/2 + lemScale*ooz*x*2)
			yp := int(height/2 - lemScale*ooz*y)
			if xp < 0 || xp >= width || yp < 0 || yp >= height {
				continue
			}

			idx := yp*width + xp
			if ooz <= l.zbuf[idx] {
				continue
			}
			l.zbuf[idx] = ooz

			rx, ry, rz := l.rotate(ox, oy, oz, sa, ca, sb, cb)
			lum := max(rx*lx+ry*ly+rz*lz, 0)
			l.buf[idx] = luminance[min(int(lum*float64(len(luminance))), len(luminance)-1)]
		}
	}

	out.WriteString("\x1b[36m")
	for y := range height {
		out.Write(l.buf[y*width : (y+1)*width])
		out.WriteByte('\n')
	}

	l.a += 0.03
	l.b += 0.05
}
