package pixeldisplay

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type PixelDisplay struct {
	scale  float64
	win    *pixelgl.Window
	imd    *imdraw.IMDraw
	pixels [64][32]uint8 // We could more efficiently use just 8 ints for the width but using a separate int per pixel keeps things relatively simple
}

func New(scale float64) *PixelDisplay {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 64*scale, 32*scale),
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	result := PixelDisplay{
		scale: scale,
		win:   win,
		imd:   imd,
	}
	return &result
}

func (pd *PixelDisplay) Closed() bool {
	return pd.win.Closed()
}

func (pd *PixelDisplay) Update() {
	// TODO: we can do this much more efficiently and just draw what's changed instead by keeping track of 'dirty' pixels
	pd.imd.Clear()
	pd.imd.Color = colornames.White

	for x, px := range pd.pixels {
		for y, enabled := range px {
			if enabled > 0 {
				pd.imd.Push(pixel.V(float64(x)*pd.scale, float64(31-y)*pd.scale), pixel.V(float64(x+1)*pd.scale, float64(31-y+1)*pd.scale))
				pd.imd.Rectangle(0)
			}
		}
	}

	pd.win.Clear(colornames.Black)
	pd.imd.Draw(pd.win)
	pd.win.Update()
}

func (pd *PixelDisplay) Draw(x uint8, y uint8, bytes []uint8) bool {
	collision := false

	for i, b := range bytes {
		for j := 0; j < 8; j++ {
			pixelSet := uint8(0)

			if (b & (0x80 >> j)) > 0 {
				pixelSet = uint8(1)
			}

			if pixelSet != 0 {
				if pd.pixels[(x+uint8(j))%64][(y+uint8(i))%32] != 0 {
					collision = true
				}
				pd.pixels[(x+uint8(j))%64][(y+uint8(i))%32] ^= pixelSet
			}
		}
	}

	return collision
}

func (pd *PixelDisplay) KeyDown(key uint8) bool {
	// Should use a map, obviously
	if key == 0x1 {
		return pd.win.Pressed(pixelgl.Key1)
	}
	if key == 0x2 {
		return pd.win.Pressed(pixelgl.Key2)
	}
	if key == 0x3 {
		return pd.win.Pressed(pixelgl.Key3)
	}
	if key == 0xC {
		return pd.win.Pressed(pixelgl.Key4)
	}
	if key == 0x4 {
		return pd.win.Pressed(pixelgl.KeyQ)
	}
	if key == 0x5 {
		return pd.win.Pressed(pixelgl.KeyW)
	}
	if key == 0x6 {
		return pd.win.Pressed(pixelgl.KeyE)
	}
	if key == 0xD {
		return pd.win.Pressed(pixelgl.KeyR)
	}
	if key == 0x7 {
		return pd.win.Pressed(pixelgl.KeyA)
	}
	if key == 0x8 {
		return pd.win.Pressed(pixelgl.KeyS)
	}
	if key == 0x9 {
		return pd.win.Pressed(pixelgl.KeyD)
	}
	if key == 0xE {
		return pd.win.Pressed(pixelgl.KeyF)
	}
	if key == 0xA {
		return pd.win.Pressed(pixelgl.KeyZ)
	}
	if key == 0x0 {
		return pd.win.Pressed(pixelgl.KeyX)
	}
	if key == 0xB {
		return pd.win.Pressed(pixelgl.KeyC)
	}
	if key == 0xF {
		return pd.win.Pressed(pixelgl.KeyV)
	}
	return false
}
