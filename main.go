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
	computer.LoadROM("test_opcode.ch8")

	for !display.Closed() && !computer.IsHalted() {
		st := time.Now()
		computer.Tick()
		display.Update(computer.GetScreen())

		tt := time.Since(st)

		time.Sleep((10 * time.Millisecond) - tt)
	}

	computer.Pause()

	// Keep the display running after halting; makes it easier to debug etc
	for !display.Closed() {
		display.Update(computer.GetScreen())
	}
}

func main() {
	pixelgl.Run(setup)
}
