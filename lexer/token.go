package lexer

type Token rune

const (
	NilToken       Token = 0
	MoveRight      Token = '>'
	MoveLeft       Token = '<'
	Increment      Token = '+'
	Decrement      Token = '-'
	Write          Token = '.'
	Read           Token = ','
	JumpIfZero     Token = '['
	JumpUnlessZero Token = ']'
)

var tokens = map[Token]struct{}{
	NilToken:       {},
	MoveRight:      {},
	MoveLeft:       {},
	Increment:      {},
	Decrement:      {},
	Write:          {},
	Read:           {},
	JumpIfZero:     {},
	JumpUnlessZero: {},
}

func NewToken(val rune) Token {
	if _, ok := tokens[Token(val)]; ok {
		return Token(val)
	}
	return NilToken
}
