package lexer

import (
	"github.com/pkg/errors"
	"io"
)

type Lexer interface {
	Lex() ([]Token, error)
}

type ByteLexer struct {
	input io.RuneReader
}

func NewLexer(input io.RuneReader) ByteLexer {
	return ByteLexer{input: input}
}

func (b ByteLexer) Lex() ([]Token, error) {
	out := make([]Token, 0)
	for {
		r, _, err := b.input.ReadRune()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return out, err
			}
			return out, nil
		}
		token := NewToken(r)
		if token != NilToken {
			out = append(out, token)
		}
	}

}
