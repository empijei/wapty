package decode

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

const eof = -1

var invalid = '�'

// Pos is an integer that define a position inside a string
type Pos int
type itemType int

func genInvalid(n int) (inv string) {
	return strings.Repeat(string(invalid), n)
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*decoder) stateFn

// decoder holds the state of the scanner.
type decoder struct {
	input string  // the string being scanned
	state stateFn // the next lexing function to enter
	pos   Pos     // current position in the input
	start Pos     // start position of this item
	width Pos     // width of last rune read from input
	out   *bytes.Buffer
}

// Constructs a new decoder
func newDecoder(input string, startState stateFn) *decoder {
	return &decoder{
		input: input,
		state: startState,
		out:   bytes.NewBuffer(nil),
	}
}

// next returns the next rune in the input.
func (l *decoder) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *decoder) peek() rune {
	p := l.width
	r := l.next()
	l.backup()
	l.width = p
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *decoder) backup() {
	l.pos -= l.width
}

// ignore skips over the pending input before this point.
func (l *decoder) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *decoder) accept(valid string) bool {
	if bytes.ContainsRune([]byte(valid), l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *decoder) acceptRun(valid string) {
	for bytes.ContainsRune([]byte(valid), l.next()) {
	}
	l.backup()
}

// decode runs the decode until EOF
func (l *decoder) decode() []byte {
	for l.state != nil {
		l.state = l.state(l)
	}
	return l.out.Bytes()
}
