package main

import (
	"math"
	"sort"
)

type RasterOptions struct {
	Wireframe bool
}

type Raster struct {
	Pixels        []byte
	W             int32
	H             int32
	BytesPerPixel int32
	Options       RasterOptions
}

// TODO: Rasterize object's vertices and draw them
func (raster Raster) DrawObject(object Mesh) {
	for i := 0; i < len(object.Vertices); i += 3 {
		v1 := Vec3to2(object.Vertices[i].Position)
		v2 := Vec3to2(object.Vertices[i+1].Position)
		v3 := Vec3to2(object.Vertices[i+2].Position)
		v1.X = (v1.X + 1) * float32(raster.W/2)
		v2.X = (v2.X + 1) * float32(raster.W/2)
		v3.X = (v3.X + 1) * float32(raster.W/2)
		v1.Y = (v1.Y + 1) * float32(raster.H/2)
		v2.Y = (v2.Y + 1) * float32(raster.H/2)
		v3.Y = (v3.Y + 1) * float32(raster.H/2)
		if raster.Options.Wireframe {
			raster.DrawLine(v1, v2)
			raster.DrawLine(v2, v3)
			raster.DrawLine(v3, v1)
		} else {
			raster.DrawTriangle(v1, v2, v3)
		}
	}
}

func (raster Raster) DrawLine(from, to Vec2) {
	if from.X > to.X {
		from.X, to.X = to.X, from.X
	}
	if from.Y > to.Y {
		from.Y, to.Y = to.Y, from.Y
	}
	//fmt.Printf("drawing from %v to %v\n", from, to)
	dx := to.X - from.X
	if dx == 0 {
		for y := from.Y; y <= to.Y; y++ {
			raster.SetPixel(int32(from.X), int32(y), 0xffffffff)
		}
		return
	}
	dy := to.Y - from.Y
	derr := math.Abs(float64(dy / dx))
	err := 0.0
	y := from.Y
	for x := from.X; x <= to.X; x++ {
		raster.SetPixel(int32(x), int32(y), 0xffffffff)
		err = err + derr
		for err >= 0.5 {
			if dy > 0 {
				y++
			} else if dy < 0 {
				y--
			}
			err = err - 1.0
		}
	}
}

func (raster Raster) SetPixel(x, y int32, pixel uint32) {
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

func (raster Raster) Clear() {
	for i := range raster.Pixels {
		raster.Pixels[i] = 0
	}
}

func (raster Raster) DrawTriangle(v1, v2, v3 Vec2) {
	// Sort ascending by Y
	verts := []Vec2{v1, v2, v3}
	sort.Slice(verts, func(i, j int) bool {
		return verts[i].Y < verts[j].Y
	})
	v1, v2, v3 = verts[0], verts[1], verts[2]

	if v2.Y == v3.Y {
		raster.drawFlatBottomTriangle(v1, v2, v3)
	} else if v1.Y == v2.Y {
		raster.drawFlagTopTriangle(v1, v2, v3)
	} else {
		v4 := Vec2{
			float32(v1.X + (v2.Y-v1.Y)/(v3.Y-v1.Y)*(v3.X-v1.X)),
			v2.Y,
		}
		raster.drawFlatBottomTriangle(v1, v2, v4)
		raster.drawFlagTopTriangle(v2, v4, v3)
		//raster.DrawLine(v1, v2)
		//raster.DrawLine(v2, v4)
		//raster.DrawLine(v4, v1)
		//raster.DrawLine(v2, v4)
		//raster.DrawLine(v4, v3)
		//raster.DrawLine(v3, v2)
	}
}

func (raster Raster) drawFlatBottomTriangle(v1, v2, v3 Vec2) {
	invslope1 := (v2.X - v1.X) / (v2.Y - v1.Y)
	invslope2 := (v3.X - v1.X) / (v3.Y - v1.Y)

	curx1 := v1.X
	curx2 := v1.X

	for scanlineY := v1.Y; scanlineY <= v2.Y; scanlineY++ {
		for x := curx1; x <= curx2; x++ {
			raster.SetPixel(int32(x), int32(scanlineY), 0xffffffff)
		}
		curx1 += invslope2
		curx2 += invslope1
	}
}

func (raster Raster) drawFlagTopTriangle(v1, v2, v3 Vec2) {
	invslope1 := (v3.X - v1.X) / (v3.Y - v1.Y)
	invslope2 := (v3.X - v2.X) / (v3.Y - v2.Y)

	curx1 := v3.X
	curx2 := v3.X

	for scanlineY := v3.Y; scanlineY > v1.Y; scanlineY-- {
		for x := curx1; x < curx2; x++ {
			raster.SetPixel(int32(x), int32(scanlineY), 0xffffffff)
		}
		curx1 -= invslope1
		curx2 -= invslope2
	}
}
