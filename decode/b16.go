package decode

import (
	"bytes"
	"encoding/hex"
)

const b16Alphabet = "0123456789abcdefABCDEF"

const b16name = "b16"

// Base16 takes a decoder and an input string
type Base16 struct {
	dec   *decoder
	input string
}

// NewB16CodecC state machine to smartly decode a string with invalid chars
// nolint: gocyclo
func NewB16CodecC(in string) CodecC {
	const (
		itemInvalid itemType = iota
		itemAlphabet
	)

	// emit should write into output what was read up until this point
	// and move l.start to l.pos
	emit := func(d *decoder, t itemType) {
		token := d.input[d.start:d.pos]

		var decodefunc func(string) []byte

		switch t {
		case itemAlphabet:
			decodefunc = func(in string) []byte {
				if len(in) < 2 {
					return []byte(genInvalid(len(in)))
				}

				odd := false
				if len(in)%2 != 0 {
					in = in[:len(in)-1]
					odd = true
				}

				buf, err := hex.DecodeString(in)
				if err != nil {
					return []byte(err.Error())
				}

				if odd {
					buf = append(buf, []byte(genInvalid(1))...)
				}
				return buf
			}

		case itemInvalid:
			decodefunc = func(in string) []byte {
				return []byte(genInvalid(len(in)))
			}
		}

		d.out.Write(decodefunc(token))
		d.start = d.pos
	}

	var (
		startState    stateFn
		invalidState  stateFn
		alphabetState stateFn
	)

	startState = func(d *decoder) stateFn {
		switch n := d.peek(); {
		case bytes.ContainsRune([]byte(b16Alphabet), n):
			return alphabetState
		case n == eof:
			return nil
		default:
			return invalidState
		}
	}

	invalidState = func(d *decoder) stateFn {
		for {
			switch n := d.next(); {
			case bytes.ContainsRune([]byte(b16Alphabet), n):
				d.backup()
				emit(d, itemInvalid)
				return alphabetState

			case n == eof:
				emit(d, itemInvalid)
				return nil
			}
		}
	}

	alphabetState = func(d *decoder) stateFn {
		for {
			switch n := d.next(); {
			case bytes.ContainsRune([]byte(b16Alphabet), n):
				d.acceptRun(b16Alphabet)
				continue

			case n == eof:
				emit(d, itemAlphabet)
				return nil

			default:
				d.backup()
				emit(d, itemAlphabet)
				return invalidState
			}
		}
	}

	return &Base16{
		dec:   newDecoder(in, startState),
		input: in,
	}
}

// Name returns the name of the codec
func (b *Base16) Name() string {
	return b16name
}

// Decode a valid b16 string
func (b *Base16) Decode() (output string) {
	return string(b.dec.decode())
}

// Encode a string into b16
func (b *Base16) Encode() (output string) {
	return hex.EncodeToString([]byte(b.input))
}

// Check returns the percentage of valid b16 characters in the input string
func (b *Base16) Check() (acceptability float64) {
	var c int
	var tot int
	for _, r := range b.input {
		tot++
		if bytes.ContainsRune([]byte(b16Alphabet), r) {
			c++
		}
	}
	//Heuristic to consider uneven strings as less likely to be valid base16
	if delta := tot % 2; delta != 0 {
		tot += delta
	}
	return float64(c) / float64(tot)
}
