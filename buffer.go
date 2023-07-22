package brainfucker

import (
	"fmt"
)

var ErrInvalidIndex = fmt.Errorf("invalid index")

type Buffer interface {
	GoRight() error
	GoLeft() error
	IncrementCurrent() error
	DecrementCurrent() error
	Current() (byte, error)
	WriteCurrent(val byte) error
	raw() []byte
}

type ByteBuffer struct {
	currentPos int
	buf        []byte
}

func NewBuffer(bufSize int) *ByteBuffer {
	return &ByteBuffer{
		currentPos: 0,
		buf:        make([]byte, bufSize),
	}
}
func (b *ByteBuffer) raw() []byte {
	return b.buf
}

func (b *ByteBuffer) GoRight() error {
	if b.currentPos+1 >= len(b.buf) {
		return ErrInvalidIndex
	}
	b.currentPos++
	return nil
}

func (b *ByteBuffer) GoLeft() error {
	if b.currentPos-1 < 0 {
		return ErrInvalidIndex
	}
	b.currentPos--
	return nil

}

func (b *ByteBuffer) IncrementCurrent() error {
	curr, err := b.Current()
	if err != nil {
		return err
	}
	return b.WriteCurrent(curr + 1)
}

func (b *ByteBuffer) DecrementCurrent() error {
	curr, err := b.Current()
	if err != nil {
		return err
	}
	return b.WriteCurrent(curr - 1)
}

func (b *ByteBuffer) Current() (byte, error) {
	if b.currentPos < 0 || b.currentPos >= len(b.buf) {
		return 0, ErrInvalidIndex
	}
	return b.buf[b.currentPos], nil
}

func (b *ByteBuffer) WriteCurrent(val byte) error {
	if b.currentPos < 0 || b.currentPos >= len(b.buf) {
		return ErrInvalidIndex
	}
	b.buf[b.currentPos] = val
	return nil
}
