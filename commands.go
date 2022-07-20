package brainfucker

type Command = rune

const (
	IncrementPointer = '>'
	DecrementPointer = '<'
	IncrementData    = '+'
	DecrementData    = '-'
	OutputData       = '.'
	InputData        = ','
	StartLoop        = '['
	EndLoop          = ']'
)
