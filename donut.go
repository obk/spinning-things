package main

import (
	"bufio"
	"math"
)

// The donut is a torus lit from a single light source; each point's
// brightness picks a character from a ramp, from dim '.' to bright '@'.

const (
	donutR1       = 1.0  // radius of the tube
	donutR2       = 2.0  // distance from the torus center to the tube center
	donutDistance = 5.0  // distance from camera to donut center
	thetaStep     = 0.07 // step around the tube's cross-section
	phiStep       = 0.02 // step around the torus's center
)

// Projection scale chosen so the donut fills about 3/4 of the screen height.
const donutScale = height * donutDistance * 3 / (8 * (donutR1 + donutR2))

const luminance = ".,-~:;=!*#$@"

type donut struct {
	a, b float64 // rotation angles around x and z
	zbuf [width * height]float64
	buf  [width * height]byte
}

func (d *donut) frame(out *bufio.Writer) {
	for i := range d.buf {
		d.buf[i] = ' '
		d.zbuf[i] = 0
	}

	sinA, cosA := math.Sincos(d.a)
	sinB, cosB := math.Sincos(d.b)

	for theta := 0.0; theta < 2*math.Pi; theta += thetaStep {
		sinT, cosT := math.Sincos(theta)
		// Point on the tube's circular cross-section, before rotation.
		circleX := donutR2 + donutR1*cosT
		circleY := donutR1 * sinT

		for phi := 0.0; phi < 2*math.Pi; phi += phiStep {
			sinP, cosP := math.Sincos(phi)

			// Spin the circle around the y axis, then rotate by A (x axis) and B (z axis).
			x := circleX*(cosB*cosP+sinA*sinB*sinP) - circleY*cosA*sinB
			y := circleX*(sinB*cosP-sinA*cosB*sinP) + circleY*cosA*cosB
			z := donutDistance + cosA*circleX*sinP + circleY*sinA
			ooz := 1 / z // "one over z": larger means closer

			// Terminal cells are roughly twice as tall as they are wide, so x is doubled.
			xp := int(width/2 + donutScale*ooz*x*2)
			yp := int(height/2 - donutScale*ooz*y)
			if xp < 0 || xp >= width || yp < 0 || yp >= height {
				continue
			}

			// Dot product of the surface normal with a light behind and above the viewer.
			lum := cosP*cosT*sinB - cosA*cosT*sinP - sinA*sinT +
				cosB*(cosA*sinT-cosT*sinA*sinP)
			if lum <= 0 {
				continue // facing away from the light
			}

			idx := yp*width + xp
			if ooz > d.zbuf[idx] {
				d.zbuf[idx] = ooz
				d.buf[idx] = luminance[min(int(lum*8), len(luminance)-1)]
			}
		}
	}

	out.WriteString("\x1b[33m")
	for y := range height {
		out.Write(d.buf[y*width : (y+1)*width])
		out.WriteByte('\n')
	}

	d.a += 0.04
	d.b += 0.02
}
