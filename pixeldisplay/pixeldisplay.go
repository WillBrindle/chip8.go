package pixeldisplay

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type PixelDisplay struct {
	scale float64
	win   *pixelgl.Window
	imd   *imdraw.IMDraw
}

func New(scale float64) *PixelDisplay {
	cfg := pixelgl.WindowConfig{
		Title:  "Chip8.go",
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

func (pd *PixelDisplay) Update(pixels *[64][32]uint8, dirty *[64][32]bool) {
	// Clear our image drawer each time. Note: Don't clear the screen
	pd.imd.Clear()

	for x, px := range *dirty {
		for y, isDirty := range px {
			// If this pixel has changed redraw it
			if isDirty {
				if pixels[x][y] > 0 {
					pd.imd.Color = colornames.White
				} else {
					pd.imd.Color = colornames.Black
				}
				pd.imd.Push(pixel.V(float64(x)*pd.scale, float64(31-y)*pd.scale), pixel.V(float64(x+1)*pd.scale, float64(31-y+1)*pd.scale))
				pd.imd.Rectangle(0)
			}
		}
	}

	pd.imd.Draw(pd.win)
	pd.win.Update()
}

func (pd *PixelDisplay) KeyDown(key uint8) bool {
	// TODO: Should use a map, obviously
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
