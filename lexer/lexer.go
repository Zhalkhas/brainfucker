package lexer

import (
	"github.com/pkg/errors"
	"io"
)

type Lexer interface {
	Lex() <-chan TokenResult
}

type ByteLexer struct {
	input io.RuneReader
}

func NewLexer(input io.RuneReader) ByteLexer {
	return ByteLexer{input: input}
}

func (b ByteLexer) Lex() <-chan TokenResult {
	out := make(chan TokenResult)
	go func() {
		defer close(out)
		for {
			r, _, err := b.input.ReadRune()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					out <- NewTokenResultError(err)
				}
				return
			}

			token := NewToken(r)
			if token != NilToken {
				out <- NewTokenResultValue(token)
			}
		}
	}()
	return out
}
