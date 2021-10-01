package chip8

type Instruction struct {
	Command   Command
	Arguments []uint16
}

type InstructionDefinition struct {
	Command   Command
	Mask      uint16
	Match     uint16
	Arguments []InstructionArgument
}

type InstructionArgument struct {
	Mask  uint16
	Shift uint16
}

var Instructions = []InstructionDefinition{
	{
		Command: CmdReturn,
		Mask:    0xFFFF,
		Match:   0x00EE,
	},
	{
		Command: CmdJump,
		Mask:    0xF000,
		Match:   0x1000,
		Arguments: []InstructionArgument{
			{Mask: 0x0fff},
		},
	},
	{
		Command: CmdCallSubRoutine,
		Mask:    0xF000,
		Match:   0x2000,
		Arguments: []InstructionArgument{
			{Mask: 0x0fff},
		},
	},
	{
		Command: CmdSkipIfEqual,
		Mask:    0xF000,
		Match:   0x3000,
		Arguments: []InstructionArgument{
			{Mask: 0x0F00, Shift: 8},
			{Mask: 0x00FF},
		},
	},
	{
		Command: CmdSkipIfNotEqual,
		Mask:    0xF000,
		Match:   0x4000,
		Arguments: []InstructionArgument{
			{Mask: 0x0F00, Shift: 8},
			{Mask: 0x00FF},
		},
	},
	{
		Command: CmdSkipIfEqualRegister,
		Mask:    0xF00F,
		Match:   0x5000,
		Arguments: []InstructionArgument{
			{Mask: 0x0F00, Shift: 8},
			{Mask: 0x00F0, Shift: 4},
		},
	},
	{
		Command: CmdSetRegister,
		Mask:    0xF000,
		Match:   0x6000,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
			{Mask: 0x00ff},
		},
	},
	{
		Command: CmdAddToRegister,
		Mask:    0xF000,
		Match:   0x7000,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
			{Mask: 0x00ff},
		},
	},
	{
		Command: CmdCopyRegister,
		Mask:    0xF00F,
		Match:   0x8000,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
			{Mask: 0x00f0, Shift: 4},
		},
	},
	{
		Command: CmdOr,
		Mask:    0xF00F,
		Match:   0x8001,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
			{Mask: 0x00f0, Shift: 4},
		},
	},
	{
		Command: CmdAnd,
		Mask:    0xF00F,
		Match:   0x8002,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
			{Mask: 0x00f0, Shift: 4},
		},
	},
	{
		Command: CmdXOr,
		Mask:    0xF00F,
		Match:   0x8003,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
			{Mask: 0x00f0, Shift: 4},
		},
	},
	{
		Command: CmdAdd,
		Mask:    0xF00F,
		Match:   0x8004,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
			{Mask: 0x00f0, Shift: 4},
		},
	},
	{
		Command: CmdSub,
		Mask:    0xF00F,
		Match:   0x8005,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
			{Mask: 0x00f0, Shift: 4},
		},
	},
	{
		Command: CmdShiftRight,
		Mask:    0xF00F,
		Match:   0x8006,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
			{Mask: 0x00f0, Shift: 4},
		},
	},
	{
		Command: CmdSubN,
		Mask:    0xF00F,
		Match:   0x8007,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
			{Mask: 0x00f0, Shift: 4},
		},
	},
	{
		Command: CmdShiftLeft,
		Mask:    0xF00F,
		Match:   0x800E,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
			{Mask: 0x00f0, Shift: 4},
		},
	},
	{
		Command: CmdSkipIfNotEqualRegister,
		Mask:    0xF00F,
		Match:   0x9000,
		Arguments: []InstructionArgument{
			{Mask: 0x0F00, Shift: 8},
			{Mask: 0x00F0, Shift: 4},
		},
	},
	{
		Command: CmdSetI,
		Mask:    0xF000,
		Match:   0xA000,
		Arguments: []InstructionArgument{
			{Mask: 0x0fff},
		},
	},
	{
		Command: CmdJumpV0Addr,
		Mask:    0xF000,
		Match:   0xB000,
		Arguments: []InstructionArgument{
			{Mask: 0x0fff},
		},
	},
	{
		Command: CmdRandom,
		Mask:    0xF000,
		Match:   0xC000,
		Arguments: []InstructionArgument{
			{Mask: 0x0F00, Shift: 8},
			{Mask: 0x00FF},
		},
	},
	{
		Command: CmdDisplaySprite,
		Mask:    0xF000,
		Match:   0xD000,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
			{Mask: 0x00f0, Shift: 4},
			{Mask: 0x000f},
		},
	},
	{
		Command: CmdSkipIfKeyPressed,
		Mask:    0xF0FF,
		Match:   0xE09E,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
		},
	},
	{
		Command: CmdSkipIfKeyNotPressed,
		Mask:    0xF0FF,
		Match:   0xE0A1,
		Arguments: []InstructionArgument{
			{Mask: 0x0f00, Shift: 8},
		},
	},
	{
		Command: CmdGetDelayTimer,
		Mask:    0xF0FF,
		Match:   0xF007,
		Arguments: []InstructionArgument{
			{Mask: 0x0F00, Shift: 8},
		},
	},
	{
		Command: CmdSetDelayTimer,
		Mask:    0xF0FF,
		Match:   0xF015,
		Arguments: []InstructionArgument{
			{Mask: 0x0F00, Shift: 8},
		},
	},
	{
		Command: CmdSetSoundTimer,
		Mask:    0xF0FF,
		Match:   0xF018,
		Arguments: []InstructionArgument{
			{Mask: 0x0F00, Shift: 8},
		},
	},
	{
		Command: CmdAddToI,
		Mask:    0xF0FF,
		Match:   0xF01E,
		Arguments: []InstructionArgument{
			{Mask: 0x0F00, Shift: 8},
		},
	},
	{
		Command: CmdSetIToFont,
		Mask:    0xF0FF,
		Match:   0xF029,
		Arguments: []InstructionArgument{
			{Mask: 0x0F00, Shift: 8},
		},
	},
	{
		Command: CmdStoreBCD,
		Mask:    0xF0FF,
		Match:   0xF033,
		Arguments: []InstructionArgument{
			{Mask: 0x0F00, Shift: 8},
		},
	},
	{
		Command: CmdReadRegisterRange,
		Mask:    0xF0FF,
		Match:   0xF055,
		Arguments: []InstructionArgument{
			{Mask: 0x0F00, Shift: 8},
		},
	},
	{
		Command: CmdReadMemoryRange,
		Mask:    0xF0FF,
		Match:   0xF065,
		Arguments: []InstructionArgument{
			{Mask: 0x0F00, Shift: 8},
		},
	},
}
