package main

import "math"

type Raster struct {
	Pixels        []byte
	W             int32
	H             int32
	BytesPerPixel int32
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
		raster.DrawLine(v1, v2)
		raster.DrawLine(v2, v3)
		raster.DrawLine(v3, v1)
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
	copy(target, []byte{0xff, 0xff, 0xff, 0xff}) // TODO
}

func (raster Raster) Clear() {
	for i := range raster.Pixels {
		raster.Pixels[i] = 0
	}
}
