package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const WIDTH = 1280
const HEIGHT = 720
const PI = 3.1415

var quit = false
var raster Raster
var camera *Camera
var paused = true
var (
	yaw   float32 = 0.0
	pitch float32 = 0.0
	roll  float32 = 0.0
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("raster", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WIDTH, HEIGHT, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
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

	//localObj := makeRGBTriangle()
	localObj := makeRGBCube()
	localObj.Transform.Position = Vec3{0, 0, 0}
	roll = 0.1

	camera = &Camera{
		Position:   Vec3{0, 0, 1},
		Up:         Vec3{0, 1, 0},
		Fov:        PI,
		Yaw:        -PI / 2,
		Pitch:      0,
		Near:       0.1,
		Far:        10,
		Width:      12.8 / 4,
		Height:     7.2 / 4,
		Projection: PROJ_ORTHO,
	}

	raster = NewRaster(WIDTH, HEIGHT)

	prevT := sdl.GetTicks()

	for {
		processInput()
		if quit {
			break
		}

		t := sdl.GetTicks()
		if !paused {
			if prevT != t {
				fmt.Printf("Loop time = %d ms (FPS = %d)\n", t-prevT, 1000/(t-prevT))
			}
			prevT = t
			//camera.Position.X = float32(math.Sin(float64(t) / 400.0))
			pitch = float32(t) / 550.0
			//pitch = yaw
			//roll = pitch
		}
		localObj.Transform.Rotation = FromEuler(yaw, pitch, roll)
		//localObj.Transform.Position.Z = 0.2*float32(math.Sin(float64(t/400.0))) - 0.5

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
		//if t%10 == 0 {
		//mx, my, _ := sdl.GetMouseState()
		//fmt.Printf("Coords: (%d, %d)\n", mx, my)
		//}

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
				case sdl.K_w:
					raster.Options.Wireframe = !raster.Options.Wireframe
				case sdl.K_p:
					paused = !paused
				case sdl.K_o:
					if camera.Projection == PROJ_ORTHO {
						camera.Projection = PROJ_PERSP
					} else {
						camera.Projection = PROJ_ORTHO
					}
				case sdl.K_d:
					raster.Options.DrawDepth = !raster.Options.DrawDepth
				case sdl.K_RIGHT:
					yaw += 0.1
				case sdl.K_LEFT:
					yaw -= 0.1
				case sdl.K_UP:
					pitch -= 0.1
				case sdl.K_DOWN:
					pitch += 0.1
				}
			}
		case *sdl.QuitEvent:
			quit = true
		}
	}
}

func makeRGBTriangle() Object {
	t := MakeObject(TRIANGLE_MESH)
	t.Vertices[0].Color = 0xff0000
	t.Vertices[1].Color = 0x00ff00
	t.Vertices[2].Color = 0x0000ff
	return t
}

func makeRGBCube() Object {
	c := MakeObject(CUBE_MESH)
	for i := 0; i < len(c.Vertices); i += 3 {
		c.Vertices[i+0].Color = 0xff0000
		c.Vertices[i+1].Color = 0x00ff00
		c.Vertices[i+2].Color = 0x0000ff
	}
	return c
}
