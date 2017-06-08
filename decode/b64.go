package decode

import (
	"bytes"
	"encoding/base64"
	"unicode/utf8"
)

type Pos int

const eof = -1

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const variant = "+/"
const padding = "="
const urlVariant = "-_"

type itemType int

const (
	itemInvalid itemType = iota
	itemAlphabet
	itemVariant
	itemUrlVariant
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

func NewB64Decoder(input string) *decoder {
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

	var decodefunc func(string) []byte

	switch t {
	case itemAlphabet, itemVariant:
		decodefunc = func(in string) []byte {
			if len(in) < 2 {
				return []byte(genInvalid(len(in)))
			}
			encoding := base64.RawStdEncoding
			buf, err := encoding.DecodeString(in)
			if err != nil {
				return []byte(err.Error())
			}
			return buf
		}

	case itemUrlVariant:
		decodefunc = func(in string) []byte {
			if len(in) < 2 {
				return []byte(genInvalid(len(in)))
			}
			encoding := base64.RawURLEncoding
			buf, err := encoding.DecodeString(in)
			if err != nil {
				return []byte(err.Error())
			}
			return buf
		}

	case itemInvalid:
		decodefunc = func(in string) []byte {
			return []byte(genInvalid(len(in)))
		}
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

func (l *decoder) decode() []byte {
	for l.state != nil {
		l.state = l.state(l)
	}
	return l.out.Bytes()
}

func startState(l *decoder) stateFn {
	switch n := l.peek(); {
	case bytes.ContainsRune([]byte(alphabet), n):
		return alphabetState
	case bytes.ContainsRune([]byte(variant), n):
		return variantState
	case bytes.ContainsRune([]byte(urlVariant), n):
		return urlVariantState
	case n == -1:
		return nil
	default:
		return invalidState
	}
}

func invalidState(l *decoder) stateFn {
	for {
		switch n := l.next(); {
		case bytes.ContainsRune([]byte(alphabet), n):
			l.backup()
			l.emit(itemInvalid)
			return alphabetState

		case bytes.ContainsRune([]byte(variant), n):
			l.backup()
			l.emit(itemInvalid)
			return variantState

		case bytes.ContainsRune([]byte(urlVariant), n):
			l.backup()
			l.emit(itemInvalid)
			return urlVariantState

		case n == -1:
			l.emit(itemInvalid)
			return nil
		}
	}
}

//TODO consider ==
func alphabetState(l *decoder) stateFn {
	for {
		switch n := l.next(); {
		case bytes.ContainsRune([]byte(alphabet), n):
			continue

		case bytes.ContainsRune([]byte(variant), n):
			l.backup()
			return variantState

		case bytes.ContainsRune([]byte(urlVariant), n):
			l.backup()
			return urlVariantState

		case n == -1:
			l.emit(itemAlphabet)
			return nil

		default:
			l.backup()
			start := l.start
			l.emit(itemAlphabet)
			ignorePadding(l, start)
			return invalidState
		}
	}
}

func variantState(l *decoder) stateFn {
	for {
		switch n := l.next(); {
		case bytes.ContainsRune([]byte(alphabet+variant), n):
			continue

		case n == -1:
			l.emit(itemVariant)
			return nil

		default:
			l.backup()
			start := l.start
			l.emit(itemVariant)
			ignorePadding(l, start)
			return invalidState
		}
	}
}

func urlVariantState(l *decoder) stateFn {
	for {
		switch n := l.next(); {
		case bytes.ContainsRune([]byte(alphabet+urlVariant), n):
			continue

		case n == -1:
			l.emit(itemUrlVariant)
			return nil

		default:
			l.backup()
			start := l.start
			l.emit(itemUrlVariant)
			ignorePadding(l, start)
			return invalidState
		}
	}
}

func ignorePadding(l *decoder, start Pos) {
	for {
		if l.peek() != '=' {
			return
		}
		switch (l.pos - start) % 4 {
		case 2, 3:
			l.next()
			l.ignore()

		default:
			return
		}
	}
}
