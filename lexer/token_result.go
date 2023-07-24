package lexer

// TokenResult represents lexing result, either token or error
type TokenResult interface {
	Token() Token
	Error() error
}

type tokenResultValue struct {
	token Token
}

func (t tokenResultValue) Token() Token {
	return t.token
}

func (t tokenResultValue) Error() error {
	return nil
}

type tokenResultError struct {
	error
}

func (t tokenResultError) Token() Token {
	return NilToken
}

func (t tokenResultError) Error() error {
	return t.error
}

func NewTokenResultValue(token Token) TokenResult {
	return tokenResultValue{token: token}
}

func NewTokenResultError(err error) TokenResult {
	return tokenResultError{err}
}
