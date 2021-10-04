package main

import (
	"chip8/chip8"
	"chip8/pixeldisplay"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func setup() {
	display := pixeldisplay.New(8)
	computer := chip8.New(display)
	computer.LoadROM("roms/pong.rom")

	ticker := time.NewTicker(time.Second / 300)

	for !display.Closed() && !computer.IsHalted() {
		computer.Tick()
		display.Update(computer.GetScreen(), computer.GetDirtyFlags())

		<-ticker.C
	}

	computer.Pause()
	ticker.Stop()

	// Keep the display running after halting; makes it easier to debug etc
	for !display.Closed() {
		display.Update(computer.GetScreen(), computer.GetDirtyFlags())
	}
}

func main() {
	pixelgl.Run(setup)
}
