package chip8

type Display interface {
	Update(*[64][32]uint8, *[64][32]bool)
	Closed() bool
	KeyDown(key uint8) bool
}
