package chip8

import (
	"math/rand"
	"testing"
)

type MockDisplay struct {
	keysDown [16]bool
}

// Stubs; we don't need these to do anything for mocking
func (md *MockDisplay) Update(*[64][32]uint8, *[64][32]bool) {}
func (md *MockDisplay) Closed() bool                         { return false }

// Returns whether a key is down or not based on storage allowing us to mock pressing buttons
func (md *MockDisplay) KeyDown(key uint8) bool { return md.keysDown[key] }

func createTestChip8(program []uint8) (*Chip8, *MockDisplay) {
	display := MockDisplay{}
	chip8 := New(&display)
	copy(chip8.memory[0x200:], program)
	return chip8, &display
}

// 2nnn/00EE - Subroutines
func Test2nnnSubroutine(t *testing.T) {
	// Our subroutine should be entered and returned from, avoiding any `0x0000` which would throw errors
	// CALL 204
	// 0x0000
	// LD V0, 2
	// RET
	// 0x0000
	chip8, _ := createTestChip8([]uint8{0x22, 0x04, 0x00, 0x00, 0x60, 0x02, 0x00, 0xEE, 0x00, 0x00})

	chip8.Tick()
	chip8.Tick()
	chip8.Tick()
	if chip8.registers[0] != 2 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 2, chip8.registers[0])
	}
}

func Test2nnnSubroutineCanCallSubRoutine(t *testing.T) {
	// Our subroutine should be entered and returned from, avoiding any 0x0000 which would throw errors
	// 200   CALL 204
	// 202   0x0000
	// 204   CALL 20A
	// 206   RET
	// 208   0x0000
	// 20A   LD V0, 2
	// 20C   RET
	chip8, _ := createTestChip8([]uint8{0x22, 0x04, 0x00, 0x00, 0x22, 0x0A, 0x00, 0xEE, 0x00, 0x00, 0x60, 0x02, 0x00, 0xEE, 0x00, 0x00})

	for i := 0; i < 5; i++ {
		chip8.Tick()
	}
	if chip8.registers[0] != 2 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 2, chip8.registers[0])
	}
}

func Test2nnnSubroutineWillThrowStackOverflows(t *testing.T) {
	// Our subroutine should be entered and returned from, avoiding any 0x0000 which would throw errors
	// 200   CALL 202
	// 202   CALL 200

	// Will be called at end to check we panicked
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	chip8, _ := createTestChip8([]uint8{0x22, 0x02, 0x22, 0x00})

	var err error = nil
	for i := 0; i < 100 && err == nil; i++ {
		err = chip8.Tick()
	}

	if err != nil {
		t.Error("Should not get an error, should panic")
	}
}

// 1nnn - jump
func Test1nnnJumpsToCorrectLocation(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x12, 0x04, 0x00, 0x00, 0x60, 0x23})

	chip8.Tick()
	err := chip8.Tick()

	if err != nil {
		t.Error("Did not skip the 0x0 instruction")
	}
	if chip8.registers[0] != 0x23 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

// 3xkk - Conditional skip
func Test3xkkJumpsIfValueEquals0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x30, 0x00, 0x00, 0x00, 0x60, 0x23})
	chip8.registers[0] = 0

	chip8.Tick()
	err := chip8.Tick()

	if err != nil {
		t.Error("Did not skip the 0x0 instruction")
	}
	if chip8.registers[0] != 0x23 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

func Test3xkkJumpsIfValueEquals255(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x30, 0xFF, 0x00, 0x00, 0x60, 0x23})
	chip8.registers[0] = 255

	chip8.Tick()
	err := chip8.Tick()

	if err != nil {
		t.Error("Did not skip the 0x0 instruction")
	}
	if chip8.registers[0] != 0x23 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

func Test3xkkDoesNotJumpIfValueNotEqual(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x30, 0xFF, 0x60, 0x23, 0x00, 0x00})
	chip8.registers[0] = 25

	chip8.Tick()
	err := chip8.Tick()

	if err != nil {
		t.Error("Did not skip the 0x0 instruction")
	}
	if chip8.registers[0] != 0x23 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

// 4xkk - not equal conditional skip
func Test4xkkDoesNotJumpIfValueEquals0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x40, 0x00, 0x60, 0x23, 0x00, 0x00})
	chip8.registers[0] = 0

	chip8.Tick()
	err := chip8.Tick()

	if err != nil {
		t.Error("Did not skip the 0x0 instruction")
	}
	if chip8.registers[0] != 0x23 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

func Test4xkkDoesNotJumpIfValueEquals255(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x40, 0xFF, 0x60, 0x23, 0x00, 0x00})
	chip8.registers[0] = 255

	chip8.Tick()
	err := chip8.Tick()

	if err != nil {
		t.Error("Did not skip the 0x0 instruction")
	}
	if chip8.registers[0] != 0x23 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

func Test4xkkDoesJumpIfValueNotEqual(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x40, 0xFF, 0x00, 0x00, 0x60, 0x23})
	chip8.registers[0] = 25

	chip8.Tick()
	err := chip8.Tick()

	if err != nil {
		t.Error("Did not skip the 0x0 instruction")
	}
	if chip8.registers[0] != 0x23 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

// 0x5xy0
func Test5xy0RegistersEqualSkipsBoth0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x50, 0x10, 0x00, 0x00, 0x60, 0x23})
	chip8.registers[0] = 0
	chip8.registers[1] = 0

	chip8.Tick()
	err := chip8.Tick()

	if err != nil {
		t.Error("Did not skip the 0x0 instruction")
	}
	if chip8.registers[0] != 0x23 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}
func Test5xy0RegistersEqualSkips(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x50, 0x10, 0x00, 0x00, 0x60, 0x23})
	chip8.registers[0] = 25
	chip8.registers[1] = 25

	chip8.Tick()
	err := chip8.Tick()

	if err != nil {
		t.Error("Did not skip the 0x0 instruction")
	}
	if chip8.registers[0] != 0x23 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

func Test5xy0RegistersNotEqualDoesNotSkip(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x50, 0x10, 0x60, 0x23, 0x00, 0x00})
	chip8.registers[0] = 25
	chip8.registers[1] = 26

	chip8.Tick()
	err := chip8.Tick()

	if err != nil {
		t.Error("Did not skip the 0x0 instruction")
	}
	if chip8.registers[0] != 0x23 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

// 0x6xkk - Set Register tests
func Test6xkkSetUninitialisedRegister(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x60, 0x23})

	chip8.Tick()

	if chip8.registers[0] != 0x23 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

func Test6xkkSetLastRegister(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x6F, 0x23})

	chip8.Tick()

	if chip8.registers[0xF] != 0x23 {
		t.Errorf("Register[0xF] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0xF])
	}
}

func Test6xkkSetAlreadyUsedRegister(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x61, 0x23})
	// Register already has a value
	chip8.registers[1] = 0x12

	chip8.Tick()

	if chip8.registers[1] != 0x23 {
		t.Errorf("Register[1] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[1])
	}
}

// 7xkk
func Test7xkkAddNumberToRegister(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x70, 23})
	chip8.registers[0] = 200

	chip8.Tick()

	if chip8.registers[0] != 223 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 223, chip8.registers[0])
	}
}

func Test7xkkAddNumberToRegisterF(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x7F, 23})
	chip8.registers[0xF] = 200

	chip8.Tick()

	if chip8.registers[0xF] != 223 {
		t.Errorf("Register[0xF] was not set correctly. Expected %d, got %d", 223, chip8.registers[0xF])
	}
}

func Test7xkkAddToRegisterWith0Value(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x70, 23})
	chip8.registers[0] = 0

	chip8.Tick()

	if chip8.registers[0] != 23 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 23, chip8.registers[0])
	}
}

func Test7xkkAdd0ToRegister(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x70, 0x00})
	chip8.registers[0] = 200

	chip8.Tick()

	if chip8.registers[0] != 200 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 200, chip8.registers[0])
	}
}

func Test7xkkAddToRegisterOverflow(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x70, 200})
	chip8.registers[0] = 200

	chip8.Tick()

	if chip8.registers[0] != 144 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 144, chip8.registers[0])
	}
}

// 8xy0
func Test8xy0SetsVxToVy(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x10})
	chip8.registers[1] = 200

	chip8.Tick()

	if chip8.registers[0] != 200 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 200, chip8.registers[0])
	}
}

func Test8xy0SetsVxToVy255(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x10})
	chip8.registers[1] = 255

	chip8.Tick()

	if chip8.registers[0] != 255 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 255, chip8.registers[0])
	}
}

func Test8xy0SetsVxToVy0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x10})
	chip8.registers[1] = 0

	chip8.Tick()

	if chip8.registers[0] != 0 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0, chip8.registers[0])
	}
}

// 8xy1
func Test8xy1Or2Registers(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x11})
	chip8.registers[0] = 0xF0
	chip8.registers[1] = 0x0F

	chip8.Tick()

	if chip8.registers[0] != 0xFF {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0xFF, chip8.registers[0])
	}
}

func Test8xy1Or2RegistersFF(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x11})
	chip8.registers[0] = 0xFF
	chip8.registers[1] = 0xFF

	chip8.Tick()

	if chip8.registers[0] != 0xFF {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0xFF, chip8.registers[0])
	}
}

func Test8xy1Or2Registers00(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x11})
	chip8.registers[0] = 0x00
	chip8.registers[1] = 0x00

	chip8.Tick()

	if chip8.registers[0] != 0x00 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x00, chip8.registers[0])
	}
}

// 8xy2
func Test8xy2And2Registers(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x12})
	chip8.registers[0] = 0xF0
	chip8.registers[1] = 0x0F

	chip8.Tick()

	if chip8.registers[0] != 0x00 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x00, chip8.registers[0])
	}
}

func Test8xy2And2RegistersFF(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x12})
	chip8.registers[0] = 0xFF
	chip8.registers[1] = 0xFF

	chip8.Tick()

	if chip8.registers[0] != 0xFF {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0xFF, chip8.registers[0])
	}
}

func Test8xy2And2Registers00(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x12})
	chip8.registers[0] = 0x00
	chip8.registers[1] = 0x00

	chip8.Tick()

	if chip8.registers[0] != 0x00 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x00, chip8.registers[0])
	}
}

// 8xy3
func Test8xy3XOr2Registers(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x13})
	chip8.registers[0] = 0xF0
	chip8.registers[1] = 0x0F

	chip8.Tick()

	if chip8.registers[0] != 0xFF {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0xFF, chip8.registers[0])
	}
}

func Test8xy3XOr2RegistersFF(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x13})
	chip8.registers[0] = 0xFF
	chip8.registers[1] = 0xFF

	chip8.Tick()

	if chip8.registers[0] != 0x00 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x00, chip8.registers[0])
	}
}

func Test8xy3XOr2Registers00(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x13})
	chip8.registers[0] = 0x00
	chip8.registers[1] = 0x00

	chip8.Tick()

	if chip8.registers[0] != 0x00 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0x00, chip8.registers[0])
	}
}

// 8xy4
func TestAdd2Registers(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x14})
	chip8.registers[0] = 23
	chip8.registers[1] = 200

	chip8.Tick()

	if chip8.registers[0] != 223 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 223, chip8.registers[0])
	}
	if chip8.registers[0xF] != 0 {
		t.Errorf("Incorrect carry byte")
	}
}

func TestAdd2RegistersWithOverflow(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x14})
	chip8.registers[0] = 200
	chip8.registers[1] = 200

	chip8.Tick()

	if chip8.registers[0] != 144 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 144, chip8.registers[0])
	}
	if chip8.registers[0xF] != 1 {
		t.Errorf("Incorrect carry byte")
	}
}

func TestAdd2RegistersWith0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x14})
	chip8.registers[0] = 23
	chip8.registers[1] = 0

	chip8.Tick()

	if chip8.registers[0] != 23 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 23, chip8.registers[0])
	}
	if chip8.registers[0xF] != 0 {
		t.Errorf("Incorrect carry byte")
	}
}

// 8xy5
func Test8xy5Subtract2Registers(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x15})
	chip8.registers[0] = 200
	chip8.registers[1] = 23

	chip8.Tick()

	if chip8.registers[0] != 177 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 177, chip8.registers[0])
	}
	if chip8.registers[0xF] != 1 {
		t.Errorf("Incorrect borrow byte")
	}
}

func Test8xy5Subtract2RegistersWithOverflow(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x15})
	chip8.registers[0] = 23
	chip8.registers[1] = 200
	chip8.Tick()

	if chip8.registers[0] != 79 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 177, chip8.registers[0])
	}
	if chip8.registers[0xF] != 0 {
		t.Errorf("Incorrect borrow byte")
	}
}

func Test8xy5Subtract0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x15})
	chip8.registers[0] = 200
	chip8.registers[1] = 0

	chip8.Tick()

	if chip8.registers[0] != 200 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 200, chip8.registers[0])
	}
	if chip8.registers[0xF] != 1 {
		t.Errorf("Incorrect borrow byte")
	}
}

// 8xy6
func Test8xy6SHREven(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x16})
	chip8.registers[0] = 200

	chip8.Tick()

	if chip8.registers[0] != 100 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 100, chip8.registers[0])
	}
	if chip8.registers[0xF] != 0 {
		t.Errorf("Incorrect lsb byte")
	}
}

func Test8xy6SHROdd(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x16})
	chip8.registers[0] = 201

	chip8.Tick()

	if chip8.registers[0] != 100 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 100, chip8.registers[0])
	}
	if chip8.registers[0xF] != 1 {
		t.Errorf("Incorrect lsb byte")
	}
}

func Test8xy6SHR0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x16})
	chip8.registers[0] = 0

	chip8.Tick()

	if chip8.registers[0] != 0 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 100, chip8.registers[0])
	}
	if chip8.registers[0xF] != 0 {
		t.Errorf("Incorrect lsb byte")
	}
}

// 8xy7
func Test8xy7Subtract2Registers(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x17})
	chip8.registers[0] = 23
	chip8.registers[1] = 200

	chip8.Tick()

	if chip8.registers[0] != 177 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 177, chip8.registers[0])
	}
	if chip8.registers[0xF] != 1 {
		t.Errorf("Incorrect borrow byte")
	}
}

func Test8xy7Subtract2RegistersWithOverflow(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x17})
	chip8.registers[0] = 200
	chip8.registers[1] = 23

	chip8.Tick()

	if chip8.registers[0] != 79 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 79, chip8.registers[0])
	}
	if chip8.registers[0xF] != 0 {
		t.Errorf("Incorrect borrow byte")
	}
}

// 8xyE
func Test8xyESHRLowNumber(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x1E})
	chip8.registers[0] = 100

	chip8.Tick()

	if chip8.registers[0] != 200 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 200, chip8.registers[0])
	}
	if chip8.registers[0xF] != 0 {
		t.Errorf("Incorrect msb byte")
	}
}

func Test8xyESHROverflow(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x1E})
	chip8.registers[0] = 200

	chip8.Tick()

	if chip8.registers[0] != 144 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 144, chip8.registers[0])
	}
	if chip8.registers[0xF] != 1 {
		t.Errorf("Incorrect msb byte")
	}
}

func Test8xyESHR0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x80, 0x1E})
	chip8.registers[0] = 0

	chip8.Tick()

	if chip8.registers[0] != 0 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0, chip8.registers[0])
	}
	if chip8.registers[0xF] != 0 {
		t.Errorf("Incorrect msb byte")
	}
}

// 9xy0
func Test9xy0SkipIfNotEqual(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x90, 0x10, 0x00, 0x00, 0x62, 0x23})
	chip8.registers[0] = 25
	chip8.registers[1] = 26

	chip8.Tick()
	err := chip8.Tick()

	if chip8.registers[2] != 0x23 {
		t.Errorf("Register[2] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[2])
	}
	if err != nil {
		t.Errorf("Hit invalid instruction")
	}
}

func Test9xy0DontSkipIfEqual(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0x90, 0x10, 0x62, 0x23})
	chip8.registers[0] = 25
	chip8.registers[1] = 25

	chip8.Tick()
	err := chip8.Tick()

	if chip8.registers[2] != 0x23 {
		t.Errorf("Register[2] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[2])
	}
	if err != nil {
		t.Errorf("Hit invalid instruction")
	}
}

// Annn
func TestAnnnSetValue(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xA0, 0x23})

	chip8.Tick()

	if chip8.memoryRegister != 0x23 {
		t.Errorf("Register[2] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[2])
	}
}

func TestAnnnSetValue0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xA0, 0x00})

	chip8.Tick()

	if chip8.memoryRegister != 0 {
		t.Errorf("Register[2] was not set correctly. Expected %d, got %d", 0, chip8.registers[2])
	}
}

// Bnnn
func TestBnnnJumpToNNNAndV0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xB0, 0x23})
	chip8.registers[0] = 25

	chip8.Tick()

	if chip8.programCounter != 60 {
		t.Errorf("PC was not set correctly. Expected %d, got %d", 60, chip8.programCounter)
	}
}

func TestBnnnJumpToNNNAndV0Both0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xB0, 0x00})
	chip8.registers[0] = 0

	chip8.Tick()

	if chip8.programCounter != 0x0 {
		t.Errorf("PC was not set correctly. Expected %d, got %d", 0x0, chip8.programCounter)
	}
}

// Cxkk
func TestCxkkRandomAnd255(t *testing.T) {
	// Seed for deterministic results
	rand.Seed(0)
	chip8, _ := createTestChip8([]uint8{0xC0, 0xFF})

	chip8.Tick()

	if chip8.registers[0] != 250 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 250, chip8.registers[0])
	}
}

func TestCxkkRandomAnd128(t *testing.T) {
	// Seed for deterministic results
	rand.Seed(0)
	chip8, _ := createTestChip8([]uint8{0xC0, 0x80})

	chip8.Tick()

	if chip8.registers[0] != 128 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 128, chip8.registers[0])
	}
}

func TestCxkkRandomAnd0(t *testing.T) {
	// Seed for deterministic results
	rand.Seed(0)
	chip8, _ := createTestChip8([]uint8{0xC0, 0x00})

	chip8.Tick()

	if chip8.registers[0] != 0 {
		t.Errorf("Register[0] was not set correctly. Expected %d, got %d", 0, chip8.registers[0])
	}
}

// Dxyn
func TestDxynDrawSingleCharacterAll1s(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xD0, 0x11})
	chip8.memory[0] = 0xFF
	chip8.registers[0] = 5
	chip8.registers[1] = 10
	chip8.registers[0xF] = 1 // set collision flag to ensure it gets reset

	// A list of pixels that we'd expect to be triggered in format of {x, y}
	expectedPixels := [][]int{{5, 10}, {6, 10}, {7, 10}, {8, 10}, {9, 10}, {10, 10}, {11, 10}, {12, 10}}

	chip8.Tick()

	pixels := chip8.GetScreen()
	dirty := chip8.GetDirtyFlags()
	for x, px := range *dirty {
		for y, isDirty := range px {
			expected := false
			for _, e := range expectedPixels {
				if e[0] == x && e[1] == y {
					expected = true
				}
			}

			if isDirty {
				if pixels[x][y] > 0 && !expected {
					t.Errorf("Pixel %d %d was unexpectedly set", x, y)
				} else if pixels[x][y] == 0 && expected {
					t.Errorf("Pixel %d %d was unexpectedly not set", x, y)
				}
			}
		}
	}

	if chip8.registers[0xF] != 0 {
		t.Errorf("Collision flag set")
	}
}

func TestDxynDrawSingleCharacterAll0s(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xD0, 0x11})
	chip8.memory[0] = 0
	chip8.registers[0] = 5
	chip8.registers[1] = 10

	// A list of pixels that we'd expect to be triggered in format of {x, y}
	expectedPixels := [][]int{}

	chip8.Tick()

	pixels := chip8.GetScreen()
	for x, px := range *pixels {
		for y := range px {
			expected := false
			for _, e := range expectedPixels {
				if e[0] == x && e[1] == y {
					expected = true
				}
			}

			if expected {
				if pixels[x][y] == 0 {
					t.Errorf("Pixel %d %d was unexpectedly not set", x, y)
				}
			} else {
				if pixels[x][y] > 0 {
					t.Errorf("Pixel %d %d was unexpectedly set", x, y)
				}
			}
		}
	}

	if chip8.registers[0xF] != 0 {
		t.Errorf("Collision flag set")
	}
}

// Overflow test

func TestDxynDrawMultipleCharacters(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xD0, 0x12})
	chip8.memory[0] = 0xFF
	chip8.memory[1] = 0x0F
	chip8.registers[0] = 5
	chip8.registers[1] = 10

	// A list of pixels that we'd expect to be triggered in format of {x, y}
	expectedPixels := [][]int{{5, 10}, {6, 10}, {7, 10}, {8, 10}, {9, 10}, {10, 10}, {11, 10}, {12, 10}, {9, 11}, {10, 11}, {11, 11}, {12, 11}}

	chip8.Tick()

	pixels := chip8.GetScreen()
	dirty := chip8.GetDirtyFlags()
	for x, px := range *dirty {
		for y, isDirty := range px {
			expected := false
			for _, e := range expectedPixels {
				if e[0] == x && e[1] == y {
					expected = true
				}
			}

			if isDirty {
				if pixels[x][y] > 0 && !expected {
					t.Errorf("Pixel %d %d was unexpectedly set", x, y)
				} else if pixels[x][y] == 0 && expected {
					t.Errorf("Pixel %d %d was unexpectedly not set", x, y)
				}
			}
		}
	}

	if chip8.registers[0xF] != 0 {
		t.Errorf("Collision flag set")
	}
}

func TestDxynDrawCollision(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xD0, 0x11, 0x60, 0x09, 0xD0, 0x11})
	chip8.memory[0] = 0xFF
	chip8.registers[0] = 5
	chip8.registers[1] = 10

	// There is an overlap so we'd expect second half to be toggled off
	expectedPixels := [][]int{{5, 10}, {6, 10}, {7, 10}, {8, 10}, {13, 10}, {14, 10}, {15, 10}, {16, 10}}

	chip8.Tick()
	chip8.Tick()
	chip8.Tick()

	pixels := chip8.GetScreen()
	dirty := chip8.GetDirtyFlags()
	for x, px := range *dirty {
		for y, isDirty := range px {
			expected := false
			for _, e := range expectedPixels {
				if e[0] == x && e[1] == y {
					expected = true
				}
			}

			if isDirty {
				if pixels[x][y] > 0 && !expected {
					t.Errorf("Pixel %d %d was unexpectedly set", x, y)
				} else if pixels[x][y] == 0 && expected {
					t.Errorf("Pixel %d %d was unexpectedly not set", x, y)
				}
			}
		}
	}

	if chip8.registers[0xF] != 1 {
		t.Errorf("Collision flag not set")
	}
}

// Ex9E
func TestEx9ECommandSkippedIfKeyPressed(t *testing.T) {
	chip8, display := createTestChip8([]uint8{0xE0, 0x9E, 0x00, 0x00, 0x60, 0x23})
	display.keysDown[0] = true

	chip8.Tick()
	chip8.Tick()

	if chip8.registers[0] != 0x23 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

func TestEx9ECommandSkippedIfKeyPressedKey2(t *testing.T) {
	chip8, display := createTestChip8([]uint8{0xE1, 0x9E, 0x00, 0x00, 0x60, 0x23})
	chip8.registers[1] = 2
	display.keysDown[2] = true

	chip8.Tick()
	chip8.Tick()

	if chip8.registers[0] != 0x23 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

func TestEx9EPanicsIfCheckingInvalidKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	chip8, display := createTestChip8([]uint8{0xE1, 0x9E, 0x00, 0x00, 0x60, 0x23})
	chip8.registers[1] = 0xFF // doesn't exist
	display.keysDown[2] = true

	chip8.Tick()
	chip8.Tick()
}

func TestEx9ECommandNotSkippedIfKeyNotPressed(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xE0, 0x9E, 0x60, 0x23, 0x00, 0x00})

	chip8.Tick()
	chip8.Tick()

	if chip8.registers[0] != 0x23 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

func TestEx9ECommandNotSkippedIfDifferentKeyPressed(t *testing.T) {
	chip8, display := createTestChip8([]uint8{0xE0, 0x9E, 0x60, 0x23, 0x00, 0x00})
	display.keysDown[1] = true

	chip8.Tick()
	chip8.Tick()

	if chip8.registers[0] != 0x23 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

// ExA1
func TestExA1CommandSkippedIfKeyNotPressed(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xE0, 0xA1, 0x00, 0x00, 0x60, 0x23})

	chip8.Tick()
	chip8.Tick()

	if chip8.registers[0] != 0x23 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

func TestExA1CommandSkippedIfKeyNotPressedKey2(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xE1, 0xA1, 0x00, 0x00, 0x60, 0x23})
	chip8.registers[1] = 2

	chip8.Tick()
	chip8.Tick()

	if chip8.registers[0] != 0x23 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

func TestExA1PanicsIfCheckingInvalidKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	chip8, display := createTestChip8([]uint8{0xE1, 0xA1, 0x00, 0x00, 0x60, 0x23})
	chip8.registers[1] = 0xFF // doesn't exist
	display.keysDown[2] = true

	chip8.Tick()
	chip8.Tick()
}

func TestExA1CommandNotSkippedIfKeyPressed(t *testing.T) {
	chip8, display := createTestChip8([]uint8{0xE0, 0xA1, 0x60, 0x23, 0x00, 0x00})
	display.keysDown[0] = true

	chip8.Tick()
	chip8.Tick()

	if chip8.registers[0] != 0x23 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

// Fx07
func TestFx15GetDelayTimer(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x07})
	chip8.delayTimer = 0x23

	chip8.Tick()

	if chip8.registers[0] != 0x23 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
}

// Fx0A
func TestFx0AWaitForKey(t *testing.T) {
	chip8, display := createTestChip8([]uint8{0xF1, 0x0A, 0x60, 0x23})

	for i := 0; i < 100; i++ {
		chip8.Tick()
	}
	// No key pressed before now so register should not be set
	if chip8.registers[0] != 0 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0, chip8.registers[0])
	}

	display.keysDown[3] = true
	chip8.Tick()
	chip8.Tick()
	if chip8.registers[0] != 0x23 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
	if chip8.registers[1] != 3 {
		t.Errorf("registers[1] was not set correctly. Expected %d, got %d", 3, chip8.registers[1])
	}
}

func TestFx0AWaitForDifferentKey(t *testing.T) {
	chip8, display := createTestChip8([]uint8{0xF1, 0x0A, 0x60, 0x23})

	for i := 0; i < 100; i++ {
		chip8.Tick()
	}
	// No key pressed before now so register should not be set
	if chip8.registers[0] != 0 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0, chip8.registers[0])
	}

	display.keysDown[8] = true
	chip8.Tick()
	chip8.Tick()
	if chip8.registers[0] != 0x23 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
	if chip8.registers[1] != 8 {
		t.Errorf("registers[1] was not set correctly. Expected %d, got %d", 8, chip8.registers[1])
	}
}

func TestFx0AWaitMultipleKeysGetsFirst(t *testing.T) {
	chip8, display := createTestChip8([]uint8{0xF1, 0x0A, 0x60, 0x23})

	for i := 0; i < 100; i++ {
		chip8.Tick()
	}
	// No key pressed before now so register should not be set
	if chip8.registers[0] != 0 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0, chip8.registers[0])
	}

	display.keysDown[8] = true
	display.keysDown[3] = true
	chip8.Tick()
	chip8.Tick()
	if chip8.registers[0] != 0x23 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 0x23, chip8.registers[0])
	}
	if chip8.registers[1] != 3 {
		t.Errorf("registers[1] was not set correctly. Expected %d, got %d", 3, chip8.registers[1])
	}
}

// Fx15
func TestFx15SetDelayTimerTo255(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x15})
	chip8.registers[0] = 255

	chip8.Tick()

	if chip8.delayTimer != 254 {
		t.Errorf("Delay Timer was not set correctly. Expected %d, got %d", 254, chip8.delayTimer)
	}
}

func TestFx15SetDelayTimerTo0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x15})
	chip8.registers[0] = 0

	chip8.Tick()

	if chip8.delayTimer != 0 {
		t.Errorf("Delay Timer was not set correctly. Expected %d, got %d", 0, chip8.delayTimer)
	}
}

// Fx18
func TestFx18SetSoundTimerTo255(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x18})
	chip8.registers[0] = 255

	chip8.Tick()

	if chip8.soundTimer != 254 {
		t.Errorf("Sound Timer was not set correctly. Expected %d, got %d", 254, chip8.soundTimer)
	}
}

func TestFx18SetSoundTimerTo0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x18})
	chip8.registers[0] = 0

	chip8.Tick()

	if chip8.soundTimer != 0 {
		t.Errorf("Sound Timer was not set correctly. Expected %d, got %d", 0, chip8.soundTimer)
	}
}

// Fx1E
func TestAddToI(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x1E})
	chip8.memoryRegister = 12
	chip8.registers[0] = 255

	chip8.Tick()

	if chip8.memoryRegister != 267 {
		t.Errorf("I was not set correctly. Expected %d, got %d", 267, chip8.memoryRegister)
	}
}

func TestAddToIOverflow(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x1E})
	chip8.memoryRegister = 0xFFFF
	chip8.registers[0] = 2

	chip8.Tick()

	if chip8.memoryRegister != 1 {
		t.Errorf("I was not set correctly. Expected %d, got %d", 1, chip8.memoryRegister)
	}
}

// Fx29
func TestFx29GoToSprite0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x29})
	chip8.registers[0] = 0

	chip8.Tick()

	if chip8.memoryRegister != 0 {
		t.Errorf("I was not set correctly. Expected %d, got %d", 0, chip8.memoryRegister)
	}
}

func TestFx29GoToSprite7(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x29})
	chip8.registers[0] = 7

	chip8.Tick()

	if chip8.memoryRegister != 35 {
		t.Errorf("I was not set correctly. Expected %d, got %d", 35, chip8.memoryRegister)
	}
}

func TestFx29GoToSpriteF(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x29})
	chip8.registers[0] = 0xF

	chip8.Tick()

	if chip8.memoryRegister != 75 {
		t.Errorf("I was not set correctly. Expected %d, got %d", 75, chip8.memoryRegister)
	}
}

func TestFx29GoToSpriteOutOfBounds(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	chip8, _ := createTestChip8([]uint8{0xF0, 0x29})
	chip8.registers[0] = 0x10

	chip8.Tick()
}

// Fx33
func TestFx33GetBCDForLessThan10(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x33})
	// Make sure these are non-zero to test they get set to 0; the value will always be less than 10 so 10 is a sutible arbitary value
	chip8.memory[0] = 10
	chip8.memory[1] = 10
	chip8.memory[2] = 10
	chip8.registers[0] = 7

	chip8.Tick()

	if chip8.memory[0] != 0 {
		t.Errorf("memory[0] was not set correctly. Expected %d, got %d", 0, chip8.memory[0])
	}
	if chip8.memory[1] != 0 {
		t.Errorf("memory[1] was not set correctly. Expected %d, got %d", 0, chip8.memory[1])
	}
	if chip8.memory[2] != 7 {
		t.Errorf("memory[2] was not set correctly. Expected %d, got %d", 7, chip8.memory[2])
	}
}

func TestFx33GetBCDForLessThan100(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x33})
	// Make sure these are non-zero to test they get set to 0; the value will always be less than 10 so 10 is a sutible arbitary value
	chip8.memory[0] = 10
	chip8.memory[1] = 10
	chip8.memory[2] = 10
	chip8.registers[0] = 29

	chip8.Tick()

	if chip8.memory[0] != 0 {
		t.Errorf("memory[0] was not set correctly. Expected %d, got %d", 0, chip8.memory[0])
	}
	if chip8.memory[1] != 2 {
		t.Errorf("memory[1] was not set correctly. Expected %d, got %d", 2, chip8.memory[1])
	}
	if chip8.memory[2] != 9 {
		t.Errorf("memory[2] was not set correctly. Expected %d, got %d", 9, chip8.memory[2])
	}
}

func TestFx33GetBCD(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x33})
	// Make sure these are non-zero to test they get set to 0; the value will always be less than 10 so 10 is a sutible arbitary value
	chip8.memory[0] = 10
	chip8.memory[1] = 10
	chip8.memory[2] = 10
	chip8.registers[0] = 143

	chip8.Tick()

	if chip8.memory[0] != 1 {
		t.Errorf("memory[0] was not set correctly. Expected %d, got %d", 1, chip8.memory[0])
	}
	if chip8.memory[1] != 4 {
		t.Errorf("memory[1] was not set correctly. Expected %d, got %d", 4, chip8.memory[1])
	}
	if chip8.memory[2] != 3 {
		t.Errorf("memory[2] was not set correctly. Expected %d, got %d", 3, chip8.memory[2])
	}
}

func TestFx33GetBCDINon0(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x33})
	// Make sure these are non-zero to test they get set to 0; the value will always be less than 10 so 10 is a sutible arbitary value
	chip8.memory[0] = 10
	chip8.memory[1] = 10
	chip8.memory[2] = 10
	chip8.memoryRegister = 100
	chip8.registers[0] = 143

	chip8.Tick()

	if chip8.memory[100] != 1 {
		t.Errorf("memory[100] was not set correctly. Expected %d, got %d", 1, chip8.memory[100])
	}
	if chip8.memory[101] != 4 {
		t.Errorf("memory[101] was not set correctly. Expected %d, got %d", 4, chip8.memory[101])
	}
	if chip8.memory[102] != 3 {
		t.Errorf("memory[102] was not set correctly. Expected %d, got %d", 3, chip8.memory[102])
	}
}

// Fx55
func TestFx55Store0RegistersInMemory(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x55})
	chip8.memory[0] = 10
	chip8.registers[0] = 7

	chip8.Tick()

	if chip8.memory[0] != 10 {
		t.Errorf("memory[0] was not set correctly. Expected %d, got %d", 10, chip8.memory[0])
	}
}

func TestFx55Store1RegistersInMemory(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF1, 0x55})
	chip8.memory[0] = 10
	chip8.memory[1] = 10
	chip8.registers[0] = 7
	chip8.registers[1] = 8

	chip8.Tick()

	if chip8.memory[0] != 7 {
		t.Errorf("memory[0] was not set correctly. Expected %d, got %d", 7, chip8.memory[0])
	}
	if chip8.memory[1] != 10 {
		t.Errorf("memory[1] was not set correctly. Expected %d, got %d", 10, chip8.memory[1])
	}
}

func TestFx55StoreAllRegistersInMemory(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xFF, 0x55})
	copy(chip8.memory[:], []uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	copy(chip8.registers[:], []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})

	chip8.Tick()

	for i := 0; i < 0xF; i++ {
		if chip8.memory[i] != uint8(i+1) {
			t.Errorf("memory[%d] was not set correctly. Expected %d, got %d", i, i+1, chip8.memory[i])
		}
	}
}

// FxFx65
func TestFx65Read0RegistersFromMemory(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF0, 0x65})
	chip8.memory[0] = 10
	chip8.registers[0] = 7

	chip8.Tick()

	if chip8.registers[0] != 7 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 7, chip8.registers[0])
	}
}

func TestFx65Read1RegistersFromMemory(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xF1, 0x65})
	chip8.memory[0] = 10
	chip8.registers[0] = 7

	chip8.Tick()

	if chip8.registers[0] != 10 {
		t.Errorf("registers[0] was not set correctly. Expected %d, got %d", 10, chip8.registers[0])
	}
}

func TestFx65ReadAllRegistersFromMemory(t *testing.T) {
	chip8, _ := createTestChip8([]uint8{0xFF, 0x65})
	copy(chip8.registers[:], []uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	copy(chip8.memory[:], []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})

	chip8.Tick()

	for i := 0; i < 0xF; i++ {
		if chip8.registers[i] != uint8(i+1) {
			t.Errorf("registers[%d] was not set correctly. Expected %d, got %d", i, i+1, chip8.registers[i])
		}
	}
}

func TestTickDecaysDelayTimer(t *testing.T) {
	// Endless loop
	chip8, _ := createTestChip8([]uint8{0x12, 0x00})
	chip8.delayTimer = 10

	chip8.Tick()
	chip8.Tick()
	chip8.Tick()

	if chip8.delayTimer != 7 {
		t.Errorf("Delay Timer was not set correctly. Expected %d, got %d", 7, chip8.delayTimer)
	}
}

func TestTickDecaysSoundTimer(t *testing.T) {
	// Endless loop
	chip8, _ := createTestChip8([]uint8{0x12, 0x00})
	chip8.soundTimer = 10

	chip8.Tick()
	chip8.Tick()
	chip8.Tick()

	if chip8.soundTimer != 7 {
		t.Errorf("Sound Timer was not set correctly. Expected %d, got %d", 7, chip8.soundTimer)
	}
}
