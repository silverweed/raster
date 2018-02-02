package main

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const WIDTH = 1280
const HEIGHT = 720
const PI = 3.1415

var quit = false

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("raster", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WIDTH, HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STATIC, WIDTH, HEIGHT)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	localObj := MakeObject(QUAD_MESH)
	localObj.Transform.Position = Vec3{0, 0, 0}

	camera := Camera{
		Position:   Vec3{0, 0, 1},
		Up:         Vec3{0, 1, 0},
		Fov:        PI,
		Yaw:        -PI / 2,
		Pitch:      0,
		Near:       0.1,
		Far:        100,
		Width:      6,
		Height:     6,
		Projection: PROJ_ORTHO,
	}

	raster := Raster{
		Pixels:        make([]byte, 4*WIDTH*HEIGHT),
		W:             WIDTH,
		H:             HEIGHT,
		BytesPerPixel: 4,
	}

	prevT := sdl.GetTicks()

	for {
		processInput()
		if quit {
			break
		}

		t := sdl.GetTicks()
		if prevT != t {
			fmt.Printf("Loop time = %d ms (FPS = %d)\n", t-prevT, 1000/(t-prevT))
		}
		prevT = t
		camera.Position.X = float32(math.Sin(float64(t) / 400.0))
		fmt.Printf("cameraX = %f\n", camera.Position.X)

		raster.Clear()

		//fmt.Printf("viewmat = %v\n", camera.ViewMatrix())

		worldObj := LocalToWorld(localObj)
		//fmt.Printf("Local: %v,\nWorld: %v\n", localObj.Mesh, worldObj)
		viewObj := WorldToView(worldObj, camera.ViewMatrix())
		//fmt.Printf("Local: %v,\nView: %v\n", localObj.Mesh, viewObj)
		clipObj := ViewToClip(viewObj, camera.ProjMatrix())
		//fmt.Printf("Local: %v,\nClip: %v\n", localObj.Mesh, clipObj)

		raster.DrawObject(clipObj)

		renderer.Clear()
		texture.Update(nil, raster.Pixels, WIDTH*4 /*4 = sizeof(uint32)*/)
		renderer.Copy(texture, nil, nil)
		renderer.Present()
		sdl.Delay(1000 / 60)
	}
}

func processInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.KeyboardEvent:
			event := event.(*sdl.KeyboardEvent)
			if event.Type == sdl.KEYDOWN {
				switch event.Keysym.Sym {
				case sdl.K_q:
					quit = true
				}
			}
		case *sdl.QuitEvent:
			quit = true
		}
	}
}
