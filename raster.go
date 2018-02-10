package main

import (
	"fmt"
	"math"
	"sort"
)

type RasterOptions struct {
	Wireframe  bool
	DrawDepth  bool
	ClearColor Color
}

type Raster struct {
	Pixels        []byte
	Depth         []float32
	W             int32
	H             int32
	BytesPerPixel int32
	Options       RasterOptions
}

func NewRaster(width, height int32) Raster {
	return Raster{
		Pixels:        make([]byte, 4*width*height),
		Depth:         make([]float32, width*height),
		W:             width,
		H:             height,
		BytesPerPixel: 4,
	}
}

// Rasterize object's vertices and draw them
func (raster Raster) DrawObject(object Mesh) {
	for i := 0; i < len(object.Vertices); i += 3 {
		v1 := Vec3to2(object.Vertices[i].Position)
		v2 := Vec3to2(object.Vertices[i+1].Position)
		v3 := Vec3to2(object.Vertices[i+2].Position)
		v1.X = (v1.X + 1) * float32(raster.W/2)
		v2.X = (v2.X + 1) * float32(raster.W/2)
		v3.X = (v3.X + 1) * float32(raster.W/2)
		v1.Y = (1 - (v1.Y+1)/2) * float32(raster.H)
		v2.Y = (1 - (v2.Y+1)/2) * float32(raster.H)
		v3.Y = (1 - (v3.Y+1)/2) * float32(raster.H)

		c1 := object.Vertices[i].Color
		c2 := object.Vertices[i+1].Color
		c3 := object.Vertices[i+2].Color
		if raster.Options.DrawDepth {
			d1 := LinearizeDepth(object.Vertices[i].Position.Z)
			d2 := LinearizeDepth(object.Vertices[i+1].Position.Z)
			d3 := LinearizeDepth(object.Vertices[i+2].Position.Z)
			c1 = ColorFromFloat(d1, d1, d1)
			c2 = ColorFromFloat(d2, d2, d2)
			c3 = ColorFromFloat(d3, d3, d3)
			fmt.Printf("z1, z2, z3 = %f, %f, %f, colors = %.6x, %.6x, %.6x\n", d1, d2, d3, c1, c2, c3)
		} else {
			// For now, if color is black (default), make it white to make the vertex visible
			for _, c := range []*Color{&c1, &c2, &c2} {
				if *c == 0 {
					*c = 0xffffff
				}
			}
		}

		//fmt.Printf("Drawing triangle from %s, %s, %s\n", v1, v2, v3)
		if raster.Options.Wireframe {
			raster.DrawLine(v1, v2, c1, c2)
			raster.DrawLine(v2, v3, c2, c3)
			raster.DrawLine(v3, v1, c3, c1)
		} else {
			z1 := object.Vertices[i].Position.Z
			z2 := object.Vertices[i+1].Position.Z
			z3 := object.Vertices[i+2].Position.Z
			raster.DrawTriangle(v1, v2, v3, c1, c2, c3, z1, z2, z3)
		}
	}
}

func (raster Raster) DrawLine(from, to Vec2, colFrom, colTo Color) {
	//fmt.Printf("Drawing line from %s to %s\n", from, to)
	dx := to.X - from.X
	if dx == 0 {
		if from.Y > to.Y {
			from.Y, to.Y = to.Y, from.Y
		}
		for y := from.Y; y <= to.Y; y++ {
			t := (y - from.Y) / (to.Y - from.Y)
			color := LerpColor(colFrom, colTo, t)
			raster.SetPixel(int32(from.X), int32(y), color)
		}
		return
	}
	dy := to.Y - from.Y
	derr := math.Abs(float64(dy / dx))
	err := 0.0
	y := from.Y
	if from.X < to.X {
		for x := from.X; x <= to.X; x++ {
			t := (x - from.X) / (to.X - from.X)
			color := LerpColor(colFrom, colTo, t)
			raster.SetPixel(int32(x), int32(y), color)
			err = err + derr
			for err >= 0.5 {
				if dy > 0 {
					y++
					if y < to.Y {
						raster.SetPixel(int32(x), int32(y), color)
					}
				} else if dy < 0 {
					y--
					if y > to.Y {
						raster.SetPixel(int32(x), int32(y), color)
					}
				}
				err = err - 1.0
			}
		}
	} else {
		for x := from.X; x >= to.X; x-- {
			t := (x - from.X) / (to.X - from.X)
			color := LerpColor(colFrom, colTo, t)
			raster.SetPixel(int32(x), int32(y), color)
			err = err + derr
			for err >= 0.5 {
				if dy > 0 {
					y++
					if y < to.Y {
						raster.SetPixel(int32(x), int32(y), color)
					}
				} else if dy < 0 {
					y--
					if y > to.Y {
						raster.SetPixel(int32(x), int32(y), color)
					}
				}
				err = err - 1.0
			}
		}
	}
}

func (raster Raster) SetPixel(x, y int32, pixel Color) {
	off := y*raster.W*4 + x*raster.BytesPerPixel
	target := raster.Pixels[off : off+raster.BytesPerPixel]
	// A
	target[3] = byte(pixel >> 24)
	// R
	target[2] = byte(pixel >> 16)
	// G
	target[1] = byte(pixel >> 8)
	// B
	target[0] = byte(pixel)
}

func (raster Raster) SetDepth(x, y int32, depth float32) {
	off := y*raster.W + x
	raster.Depth[off] = depth
}

func (raster Raster) GetDepth(x, y int32) float32 {
	if x >= raster.W || y >= raster.H {
		panic(fmt.Sprintf("x or y out of bounds! (x = %d / %d, y = %d / %d)", x, raster.W, y, raster.H))
	}
	return raster.Depth[y*raster.W+x]
}

func (raster Raster) Clear() {
	for i := 0; i < len(raster.Pixels)/4; i++ {
		raster.Pixels[4*i] = byte(raster.Options.ClearColor >> 24)
		raster.Pixels[4*i+1] = byte(raster.Options.ClearColor >> 16)
		raster.Pixels[4*i+2] = byte(raster.Options.ClearColor >> 8)
		raster.Pixels[4*i+2] = byte(raster.Options.ClearColor)
		raster.Depth[i] = 0
	}
}

func (raster Raster) DrawTriangle(v1, v2, v3 Vec2, c1, c2, c3 Color, z1, z2, z3 float32) {
	// Sort ascending by Y
	verts := []Vec2{v1, v2, v3}
	sort.Slice(verts, func(i, j int) bool {
		return verts[i].Y < verts[j].Y
	})
	v1, v2, v3 = verts[0], verts[1], verts[2]

	if v2.Y == v3.Y {
		raster.drawFlatBottomTriangle(v1, v2, v3, c1, c2, c3, z1, z2, z3)
	} else if v1.Y == v2.Y {
		raster.drawFlagTopTriangle(v1, v2, v3, c1, c2, c3, z1, z2, z3)
	} else {
		v4 := Vec2{
			float32(v1.X + (v2.Y-v1.Y)/(v3.Y-v1.Y)*(v3.X-v1.X)),
			v2.Y,
		}
		z4 := LerpF32Barycentric(v4, v1, v2, v3, z1, z2, z3)
		c4 := LerpColorBarycentric(v4, v1, v2, v3, c1, c2, c3)
		raster.drawFlatBottomTriangle(v1, v2, v4, c1, c2, c4, z1, z2, z4)
		raster.drawFlagTopTriangle(v2, v4, v3, c2, c4, c3, z2, z4, z3)
	}
}

func (raster Raster) drawFlatBottomTriangle(v1, v2, v3 Vec2, c1, c2, c3 Color, z1, z2, z3 float32) {
	invslope1 := (v2.X - v1.X) / (v2.Y - v1.Y)
	invslope2 := (v3.X - v1.X) / (v3.Y - v1.Y)

	curx1 := v1.X
	curx2 := v1.X

	for scanlineY := v1.Y; scanlineY <= v2.Y; scanlineY++ {
		if curx1 < curx2 {
			for x := curx1; x <= curx2; x++ {
				p := Vec2{x, scanlineY}
				color := LerpColorBarycentric(p, v1, v2, v3, c1, c2, c3)
				z := LerpF32Barycentric(p, v1, v2, v3, z1, z2, z3)
				x, scanlineY := int32(x), int32(scanlineY)
				if z < raster.GetDepth(x, scanlineY) {
					raster.SetPixel(x, scanlineY, color)
					raster.SetDepth(x, scanlineY, z)
				}
			}
		} else {
			for x := curx2; x <= curx1; x++ {
				p := Vec2{x, scanlineY}
				color := LerpColorBarycentric(p, v1, v2, v3, c1, c2, c3)
				z := LerpF32Barycentric(p, v1, v2, v3, z1, z2, z3)
				x, scanlineY := int32(x), int32(scanlineY)
				if z < raster.GetDepth(x, scanlineY) {
					raster.SetPixel(x, scanlineY, color)
					raster.SetDepth(x, scanlineY, z)
				}
			}

		}
		curx1 += invslope2
		curx2 += invslope1
	}
}

func (raster Raster) drawFlagTopTriangle(v1, v2, v3 Vec2, c1, c2, c3 Color, z1, z2, z3 float32) {
	invslope1 := (v3.X - v1.X) / (v3.Y - v1.Y)
	invslope2 := (v3.X - v2.X) / (v3.Y - v2.Y)

	curx1 := v3.X
	curx2 := v3.X

	for scanlineY := v3.Y; scanlineY >= v1.Y; scanlineY-- {
		if curx1 < curx2 {
			for x := curx1; x <= curx2; x++ {
				p := Vec2{x, scanlineY}
				color := LerpColorBarycentric(p, v1, v2, v3, c1, c2, c3)
				z := LerpF32Barycentric(p, v1, v2, v3, z1, z2, z3)
				x, scanlineY := int32(x), int32(scanlineY)
				if z < raster.GetDepth(x, scanlineY) {
					raster.SetPixel(x, scanlineY, color)
					raster.SetDepth(x, scanlineY, z)
				}
			}
		} else {
			for x := curx2; x <= curx1; x++ {
				p := Vec2{x, scanlineY}
				color := LerpColorBarycentric(p, v1, v2, v3, c1, c2, c3)
				z := LerpF32Barycentric(p, v1, v2, v3, z1, z2, z3)
				x, scanlineY := int32(x), int32(scanlineY)
				if z < raster.GetDepth(x, scanlineY) {
					raster.SetPixel(x, scanlineY, color)
					raster.SetDepth(x, scanlineY, z)
				}
			}

		}
		curx1 -= invslope1
		curx2 -= invslope2
	}
}

func LerpColor(startcol, endcol Color, t float32) Color {
	lerpCol := func(from, to Color, t float32) Color {
		return Color((1.0-t)*float32(from) + t*float32(to))
	}
	startred, endred := (startcol>>16)&0xff, (endcol>>16)&0xff
	startgreen, endgreen := (startcol>>8)&0xff, (endcol>>8)&0xff
	startblue, endblue := startcol&0xff, endcol&0xff
	return 0xFF<<24 | lerpCol(startred, endred, t)<<16 | lerpCol(startgreen, endgreen, t)<<8 | lerpCol(startblue, endblue, t)
}

func LerpColorBarycentric(p, v1, v2, v3 Vec2, col1, col2, col3 Color) Color {
	a, b, c := Barycentric(p, v1, v2, v3)
	red := Color(float32((col1>>16)&0xff)*a + float32((col2>>16)&0xff)*b + float32((col3>>16)&0xff)*c)
	green := Color(float32((col1>>8)&0xff)*a + float32((col2>>8)&0xff)*b + float32((col3>>8)&0xff)*c)
	blue := Color(float32((col1)&0xff)*a + float32((col2)&0xff)*b + float32((col3)&0xff)*c)
	return 0xff<<24 | red<<16 | green<<8 | blue
}

func LerpF32(from, to, t float32) float32 {
	return (1.0-t)*from + t*to
}

func LerpF32Barycentric(p, v1, v2, v3 Vec2, f1, f2, f3 float32) float32 {
	a, b, c := Barycentric(p, v1, v2, v3)
	return f1*a + f2*b + f3*c
}

func Barycentric(p, va, vb, vc Vec2) (a float32, b float32, c float32) {
	v0 := vb.Add(va.Neg())
	v1 := vc.Add(va.Neg())
	v2 := p.Add(va.Neg())
	d00 := Dot2(v0, v0)
	d01 := Dot2(v0, v1)
	d11 := Dot2(v1, v1)
	d20 := Dot2(v2, v0)
	d21 := Dot2(v2, v1)
	denom := d00*d11 - d01*d01
	b = (d11*d20 - d01*d21) / denom
	c = (d00*d21 - d01*d20) / denom
	a = 1.0 - b - c
	return
}

func LinearizeDepth(depth float32) float32 {
	return 1 - 2.5*(1+depth)
	//return (2 * camera.Near) / (camera.Far + camera.Near - depth*(camera.Far-camera.Near))
}

func ColorFromFloat(r, g, b float32) Color {
	return Color(0xff<<24 | uint32(r*255)<<16 | uint32(g*255)<<8 | uint32(b*255))
}
