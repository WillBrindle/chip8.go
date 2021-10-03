package chip8

type Command int16

const (
	CmdUndefined Command = iota
	CmdAddToRegister
	CmdAdd
	CmdAddToI
	CmdAnd
	CmdCall
	CmdCallSubRoutine
	CmdCopyRegister
	CmdClear
	CmdDisplaySprite
	CmdGetDelayTimer
	CmdJump
	CmdJumpV0Addr
	CmdOr
	CmdRandom
	CmdReadMemoryRange
	CmdReadRegisterRange
	CmdReturn
	CmdSetDelayTimer
	CmdSetI
	CmdSetIToFont
	CmdSetRegister
	CmdSetSoundTimer
	CmdShiftLeft
	CmdShiftRight
	CmdSkipIfEqual
	CmdSkipIfEqualRegister
	CmdSkipIfKeyPressed
	CmdSkipIfNotEqual
	CmdSkipIfNotEqualRegister
	CmdSkipIfKeyNotPressed
	CmdStoreBCD
	CmdSub
	CmdSubN
	CmdWaitForKey
	CmdXOr
)
