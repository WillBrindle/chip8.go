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

// Keypad:
// 1 2 3 4
// Q W E R
// A S D F
// Z X C V
var Keys = map[uint8]pixelgl.Button{
	0x1: pixelgl.Key1,
	0x2: pixelgl.Key2,
	0x3: pixelgl.Key3,
	0xC: pixelgl.Key4,
	0x4: pixelgl.KeyQ,
	0x5: pixelgl.KeyW,
	0x6: pixelgl.KeyE,
	0xD: pixelgl.KeyR,
	0x7: pixelgl.KeyA,
	0x8: pixelgl.KeyS,
	0x9: pixelgl.KeyD,
	0xE: pixelgl.KeyF,
	0xA: pixelgl.KeyZ,
	0x0: pixelgl.KeyX,
	0xB: pixelgl.KeyC,
	0xF: pixelgl.KeyV,
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
	// Clear our image drawer each time. Note: Don't clear the screen. This means we're only drawing the changes each time rather than
	// redrawing everything.
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
	return pd.win.Pressed(Keys[key])
}
