package interpreter

type Tape struct {
	buf  []byte
	curr int
}

func NewTape(bufSize int) Tape {
	return Tape{
		buf:  make([]byte, 0, bufSize),
		curr: 0,
	}
}

func (s *Tape) Curr() byte {
	if len(s.buf) <= s.curr {
		s.buf = s.buf[:s.curr+1]
	}
	return s.buf[s.curr]
}

func (s *Tape) WriteCurr(val byte) {
	if len(s.buf) <= s.curr {
		s.buf = s.buf[:s.curr+1]
	}
	s.buf[s.curr] = val
}

func (s *Tape) GoRight() {
	if s.curr < cap(s.buf) {
		s.curr++
	}
}
func (s *Tape) GoLeft() {
	if s.curr > 0 {
		s.curr--
	}
}
