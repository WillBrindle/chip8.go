package chip8

import (
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
)

var Fonts = []uint8{
	// 0
	0xF0, 0x90, 0x90, 0x90, 0xF0,
	// 1
	0x20, 0x60, 0x20, 0x20, 0x70,
	// 2
	0xF0, 0x10, 0xF0, 0x80, 0xF0,
	// 3
	0xF0, 0x10, 0xF0, 0x10, 0xF0,
	// 4
	0x90, 0x90, 0xF0, 0x10, 0x10,
	// 5
	0xF0, 0x80, 0xF0, 0x10, 0xF0,
	// 6
	0xF0, 0x80, 0xF0, 0x90, 0xF0,
	// 7
	0xF0, 0x10, 0x20, 0x40, 0x40,
	// 8
	0xF0, 0x90, 0xF0, 0x90, 0xF0,
	// 9
	0xF0, 0x90, 0xF0, 0x10, 0xF0,
	// A
	0xF0, 0x90, 0xF0, 0x90, 0x90,
	// B
	0xE0, 0x90, 0xE0, 0x90, 0xE0,
	// C
	0xF0, 0x80, 0x80, 0x80, 0xF0,
	// D,
	0xE0, 0x90, 0x90, 0x90, 0xE0,
	// E
	0xF0, 0x80, 0xF0, 0x80, 0xF0,
	// F
	0xF0, 0x80, 0xF0, 0x80, 0x80,
}

type Chip8 struct {
	display        Display
	memory         [4096]uint8
	registers      [16]uint8
	stack          [16]uint16
	screen         [64][32]uint8 // We could more efficiently use just 8 ints for the width but using a separate int per pixel keeps things relatively simple
	memoryRegister uint16        // Refered as 'I' in documentation
	programCounter uint16
	stackPointer   int8
	delayTimer     uint8
	soundTimer     uint8
	halted         bool
}

func New(display Display) *Chip8 {
	// put fonts into memory
	memory := [4096]uint8{}
	copy(memory[0:], Fonts[:])

	chip8 := Chip8{
		display:        display,
		halted:         false,
		programCounter: 0x200,
		stackPointer:   -1,
		memory:         memory,
	}

	return &chip8
}

func parseArguments(arguments []InstructionArgument, val uint16) []uint16 {
	result := []uint16{}

	for _, v := range arguments {
		result = append(result, (val&v.Mask)>>v.Shift)
	}

	return result
}

func parseInstruction(val uint16) (*Instruction, error) {
	for _, v := range Instructions {
		if val&v.Mask == v.Match {
			instr := Instruction{
				Command:   v.Command,
				Arguments: parseArguments(v.Arguments, val),
			}
			return &instr, nil
		}
	}

	return nil, errors.New("Unhandled Instruction 0x" + strconv.FormatInt(int64(val), 16))
}

func (c8 *Chip8) LoadROM(rom string) error {
	bytes, err := ioutil.ReadFile(rom)
	if err != nil {
		return err
	}

	copy(c8.memory[0x200:], bytes[:])
	return nil
}

func (c8 *Chip8) LoadFromMemory(data []uint8) error {
	copy(c8.memory[0x200:], data[:])
	return nil
}

func (c8 *Chip8) next() {
	c8.programCounter += 2
}

func (c8 *Chip8) readInstruction() (*Instruction, error) {
	val := (uint16(c8.memory[c8.programCounter]) << 8) | uint16(c8.memory[c8.programCounter+1])
	return parseInstruction(val)
}

func (c8 *Chip8) setRegister(register uint16, value uint8) {
	c8.registers[register] = value
}

func (c8 *Chip8) setI(value uint16) {
	c8.memoryRegister = value
}

func (c8 *Chip8) jump(position uint16) {
	c8.programCounter = position
}

func (c8 *Chip8) jumpV0AndAddr(addr uint16) {
	c8.programCounter = uint16(c8.registers[0]) + addr
}

func (c8 *Chip8) drawSprite(register1 uint16, register2 uint16, nibble uint16) {
	x := c8.registers[register1]
	y := c8.registers[register2]
	bytes := c8.memory[c8.memoryRegister : c8.memoryRegister+nibble]

	// Reset the collision flag to 0
	c8.registers[0xF] = 0

	for i, b := range bytes {
		for j := 0; j < 8; j++ {
			pixelSet := uint8(0)

			if (b & (0x80 >> j)) > 0 {
				pixelSet = uint8(1)
			}

			if pixelSet != 0 {
				// A collision occured so set to 1
				if c8.screen[(x+uint8(j))%64][(y+uint8(i))%32] != 0 {
					c8.registers[0xF] = 1
				}
				c8.screen[(x+uint8(j))%64][(y+uint8(i))%32] ^= pixelSet
			}
		}
	}
}

func (c8 *Chip8) callSubroutine(address uint16) {
	c8.stackPointer += 1
	if c8.stackPointer > 15 {
		panic("Stack overflow")
	}
	c8.stack[c8.stackPointer] = c8.programCounter
	c8.programCounter = address
}

func (c8 *Chip8) storeBCD(register uint16) {
	value := c8.registers[register]
	c8.memory[c8.memoryRegister] = uint8((value / 100) % 10)
	c8.memory[c8.memoryRegister+1] = uint8((value / 10) % 10)
	c8.memory[c8.memoryRegister+2] = uint8(value % 10)
}

func (c8 *Chip8) readMemoryRange(num uint16) {
	for i := 0; i < int(num); i++ {
		c8.registers[i] = c8.memory[int(c8.memoryRegister)+i]
	}
}

func (c8 *Chip8) readRegisterRange(num uint16) {
	for i := 0; i < int(num); i++ {
		c8.memory[int(c8.memoryRegister)+i] = c8.registers[i]
	}
}

func (c8 *Chip8) setLocationToFont(register uint16) {
	c8.memoryRegister = uint16(c8.registers[register]) * 5
}

func (c8 *Chip8) addToRegister(register uint16, value uint8) {
	c8.registers[register] += value
}

func (c8 *Chip8) returnFromSubroutine() {
	c8.programCounter = c8.stack[c8.stackPointer]
	c8.stackPointer--
}

func (c8 *Chip8) setDelayTimer(register uint16) {
	c8.delayTimer = c8.registers[register]
}

func (c8 *Chip8) setSoundTimer(register uint16) {
	c8.soundTimer = c8.registers[register]
}

func (c8 *Chip8) putDelayTimerIntoRegister(register uint16) {
	c8.registers[register] = c8.delayTimer
}

func (c8 *Chip8) skipIfEqual(register uint16, value uint8) {
	if c8.registers[register] == value {
		c8.programCounter += 2
	}
}

func (c8 *Chip8) skipIfNotEqual(register uint16, value uint8) {
	if c8.registers[register] != value {
		c8.programCounter += 2
	}
}

func (c8 *Chip8) skipIfEqualRegister(register1 uint16, register2 uint16) {
	if c8.registers[register1] == c8.registers[register2] {
		c8.programCounter += 2
	}
}

func (c8 *Chip8) skipIfNotEqualRegister(register1 uint16, register2 uint16) {
	if c8.registers[register1] != c8.registers[register2] {
		c8.programCounter += 2
	}
}

func (c8 *Chip8) random(register uint16, value uint8) {
	random := uint8(rand.Intn(256))
	c8.registers[register] = random & value
}

func (c8 *Chip8) skipIfKeyPressed(register uint16, keyPressed bool) {
	if c8.display.KeyDown(c8.registers[register]) == keyPressed {
		c8.programCounter += 2
	}
}

func (c8 *Chip8) and(register1 uint16, register2 uint16) {
	c8.registers[register1] &= c8.registers[register2]
}

func (c8 *Chip8) or(register1 uint16, register2 uint16) {
	c8.registers[register1] |= c8.registers[register2]
}

func (c8 *Chip8) xor(register1 uint16, register2 uint16) {
	c8.registers[register1] ^= c8.registers[register2]
}

func (c8 *Chip8) shiftRight(register1 uint16) {
	c8.registers[0xF] = c8.registers[register1] & 0x1
	c8.registers[register1] /= 2
}

func (c8 *Chip8) shiftLeft(register1 uint16) {
	c8.registers[0xF] = (c8.registers[register1] & 0x80) >> 7
	c8.registers[register1] *= 2
}

func (c8 *Chip8) add(register1 uint16, register2 uint16) {
	vx := uint16(c8.registers[register1])
	vy := uint16(c8.registers[register2])

	res := vx + vy
	if res > 255 {
		c8.registers[0xF] = 1
	} else {
		c8.registers[0xF] = 0
	}

	c8.registers[register1] = uint8(res & 0xFF)
}

func (c8 *Chip8) sub(register1 uint16, register2 uint16) {
	vx := c8.registers[register1]
	vy := c8.registers[register2]

	if vx > vy {
		c8.registers[0xF] = 1
	} else {
		c8.registers[0xF] = 0
	}

	c8.registers[register1] = vx - vy
}

func (c8 *Chip8) subN(register1 uint16, register2 uint16) {
	vx := c8.registers[register1]
	vy := c8.registers[register2]

	if vy > vx {
		c8.registers[0xF] = 1
	} else {
		c8.registers[0xF] = 0
	}

	c8.registers[register1] = vy - vx
}

func (c8 *Chip8) addToI(register uint16) {
	c8.memoryRegister += uint16(c8.registers[register])
}

func (c8 *Chip8) copyRegister(register1 uint16, register2 uint16) {
	c8.registers[register1] = c8.registers[register2]
}

func (c8 *Chip8) Tick() error {
	instruction, err := c8.readInstruction()

	// TODO: in future we'll want this to bubble up but exiting here is fine for now
	if err != nil {
		log.Println(err)
		c8.Pause()
		return err
	}

	goToNextInstruction := true

	switch instruction.Command {
	case CmdSetRegister:
		c8.setRegister(instruction.Arguments[0], uint8(instruction.Arguments[1]))
	case CmdSetI:
		c8.setI(instruction.Arguments[0])
	case CmdDisplaySprite:
		c8.drawSprite(instruction.Arguments[0], instruction.Arguments[1], instruction.Arguments[2])
	case CmdCallSubRoutine:
		c8.callSubroutine(instruction.Arguments[0])
		goToNextInstruction = false
	case CmdStoreBCD:
		c8.storeBCD(instruction.Arguments[0])
	case CmdReadMemoryRange:
		c8.readMemoryRange(instruction.Arguments[0])
	case CmdReadRegisterRange:
		c8.readRegisterRange(instruction.Arguments[0])
	case CmdSetIToFont:
		c8.setLocationToFont(instruction.Arguments[0])
	case CmdAddToRegister:
		c8.addToRegister(instruction.Arguments[0], uint8(instruction.Arguments[1]))
	case CmdReturn:
		c8.returnFromSubroutine()
	case CmdSetDelayTimer:
		c8.setDelayTimer(instruction.Arguments[0])
	case CmdSetSoundTimer:
		c8.setSoundTimer(instruction.Arguments[0])
	case CmdGetDelayTimer:
		c8.putDelayTimerIntoRegister(instruction.Arguments[0])
	case CmdSkipIfEqual:
		c8.skipIfEqual(instruction.Arguments[0], uint8(instruction.Arguments[1]))
	case CmdSkipIfNotEqual:
		c8.skipIfNotEqual(instruction.Arguments[0], uint8(instruction.Arguments[1]))
	case CmdSkipIfEqualRegister:
		c8.skipIfEqualRegister(instruction.Arguments[0], instruction.Arguments[1])
	case CmdSkipIfNotEqualRegister:
		c8.skipIfNotEqualRegister(instruction.Arguments[0], instruction.Arguments[1])
	case CmdJump:
		c8.jump(instruction.Arguments[0])
		goToNextInstruction = false
	case CmdRandom:
		c8.random(instruction.Arguments[0], uint8(instruction.Arguments[1]))
	case CmdSkipIfKeyNotPressed:
		c8.skipIfKeyPressed(instruction.Arguments[0], false)
	case CmdSkipIfKeyPressed:
		c8.skipIfKeyPressed(instruction.Arguments[0], true)
	case CmdAnd:
		c8.and(instruction.Arguments[0], instruction.Arguments[1])
	case CmdOr:
		c8.or(instruction.Arguments[0], instruction.Arguments[1])
	case CmdXOr:
		c8.xor(instruction.Arguments[0], instruction.Arguments[1])
	case CmdAdd:
		c8.add(instruction.Arguments[0], instruction.Arguments[1])
	case CmdSub:
		c8.sub(instruction.Arguments[0], instruction.Arguments[1])
	case CmdSubN:
		c8.subN(instruction.Arguments[0], instruction.Arguments[1])
	case CmdCopyRegister:
		c8.copyRegister(instruction.Arguments[0], instruction.Arguments[1])
	case CmdShiftRight:
		c8.shiftRight(instruction.Arguments[0])
	case CmdShiftLeft:
		c8.shiftLeft(instruction.Arguments[0])
	case CmdJumpV0Addr:
		c8.jumpV0AndAddr(instruction.Arguments[0])
		goToNextInstruction = false
	case CmdAddToI:
		c8.addToI(instruction.Arguments[0])
	}

	if c8.delayTimer > 0 {
		c8.delayTimer--
	}
	if c8.soundTimer > 0 {
		c8.soundTimer--
	}

	if goToNextInstruction {
		c8.next()
	}

	return nil
}

func (c8 *Chip8) IsHalted() bool {
	return c8.halted
}

func (c8 *Chip8) Pause() {
	c8.halted = true
}

func (c8 *Chip8) GetScreen() *[64][32]uint8 {
	return &c8.screen
}
