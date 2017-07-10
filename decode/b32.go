package decode

import (
	"bytes"
	"encoding/base32"
	"strings"
)

// TODO add state that handels = as padding and invalid chars

const b32Alphabet = "abcdefghijklmnopqrstuvwxyz234567ABCDEFGHIJKLMNOPQRSTUVWXYZ="

const b32name = "b32"

// Base32 takes a decoder and an input string
type Base32 struct {
	dec   *decoder
	input string
}

// NewB32CodecC state machine to smartly decode a string with invalid chars
// nolint: gocyclo
func NewB32CodecC(in string) CodecC {
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

				in = strings.ToUpper(in)
				odd := false

				// checking if len(in) is correct, then add padding
				switch n := len(in) % 8; n {
				case 6, 3, 1:
					in = in[:len(in)-1]
					odd = true
				}

				if len(in)%8 != 0 {
					pad := 8 - len(in)%8
					in = in + strings.Repeat("=", pad)
				}

				encoding := base32.StdEncoding
				buf, err := encoding.DecodeString(in)
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
		case bytes.ContainsRune([]byte(b32Alphabet), n):
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
			case bytes.ContainsRune([]byte(b32Alphabet), n):
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
			case bytes.ContainsRune([]byte(b32Alphabet), n):
				d.acceptRun(b32Alphabet)
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

	return &Base32{
		dec:   newDecoder(in, startState),
		input: in,
	}
}

// Name returns the name of the codec
func (b *Base32) Name() string {
	return b32name
}

// Decode a valid b32 string
func (b *Base32) Decode() (output string) {
	return string(b.dec.decode())
}

// Encode a string into b32
func (b *Base32) Encode() (output string) {
	return base32.StdEncoding.EncodeToString([]byte(b.input))
}

// Check returns the percentage of valid b32 characters in the input string
func (b *Base32) Check() (acceptability float64) {
	var c int
	var tot int
	for _, r := range b.input {
		tot++
		if bytes.ContainsRune([]byte(b32Alphabet), r) {
			c++
		}
	}
	//Heuristic to consider uneven strings as less likely to be valid base32
	if delta := tot % 2; delta != 0 {
		tot += delta
	}
	return float64(c) / float64(tot)
}
