package parse

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

type Pos int

const eof = -1

type itemType int

const (
	itemInvalid itemType = iota
)

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

func NewLexer(input string) *decoder {
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

// emit should write into output what was read up until this point and move l.start to l.pos
func (l *decoder) emit(t itemType) {
	token := l.input[l.start:l.pos]
	//TODO switch on itemtype to select proper decode function

	decodefunc := func(s string) []byte {
		//This is a null decoder, implement this!
		return []byte(s)
	}

	l.out.Write(decodefunc(token))
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *decoder) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *decoder) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *decoder) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *decoder) decode() []byte {
	for l.state != nil {
		l.state = l.state(l)
	}
	return l.out.Bytes()
}

func startState(l *decoder) stateFn {
	//TODO
	panic("Not implemented yet")
}

//TODO other states
