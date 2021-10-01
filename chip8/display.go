package chip8

type Display interface {
	Update()
	Closed() bool
	Draw(x uint8, y uint8, bytes []uint8) bool
	KeyDown(key uint8) bool
}
