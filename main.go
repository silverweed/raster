package main

import (
	"fmt"

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

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

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

	for {
		processInput()
		if quit {
			break
		}

		worldObj := LocalToWorld(localObj)
		//fmt.Printf("Local: %v,\nWorld: %v\n", localObj.Mesh, worldObj)
		viewObj := WorldToView(worldObj, camera.ViewMatrix())
		//fmt.Printf("Local: %v,\nView: %v\n", localObj.Mesh, viewObj)
		clipObj := ViewToClip(viewObj, camera.ProjMatrix())
		fmt.Printf("Local: %v,\nClip: %v\n", localObj.Mesh, clipObj)

		DrawObject(surface, clipObj)

		//surface.FillRect(&rect, 0xffff0000)
		window.UpdateSurface()
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
