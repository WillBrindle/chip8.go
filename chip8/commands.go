package chip8

type Command int16

const (
	CmdUndefined Command = iota
	CmdClear
	CmdReturn
	CmdJump
	CmdCall
	CmdSetRegister
	CmdSetI
	CmdDisplaySprite
	CmdCallSubRoutine
	CmdStoreBCD
	CmdReadMemoryRange
	CmdReadRegisterRange
	CmdSetIToFont
	CmdAddToRegister
	CmdSetDelayTimer
	CmdSetSoundTimer
	CmdGetDelayTimer
	CmdSkipIfEqual
	CmdSkipIfNotEqual
	CmdSkipIfEqualRegister
	CmdSkipIfNotEqualRegister
	CmdRandom
	CmdSkipIfKeyPressed
	CmdSkipIfKeyNotPressed
	CmdAnd
	CmdOr
	CmdXOr
	CmdAdd
	CmdSub
	CmdSubN
	CmdCopyRegister
	CmdShiftRight
	CmdShiftLeft
	CmdJumpV0Addr
	CmdAddToI
)
