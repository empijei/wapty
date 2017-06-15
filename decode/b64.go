package decode

import (
	"bytes"
	"encoding/base64"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const variant = "+/"
const padding = "="
const urlVariant = "-_"

const b64name = "b64"

const (
	itemInvalid itemType = iota
	itemAlphabet
	itemVariant
	itemUrlVariant
)

func init() {
	addCodecC(b64name, codecConstructor(NewB64CodecC))
}

type Base64 struct {
	dec   *decoder
	input string
}

// nolint: gocyclo
func NewB64CodecC(in string) CodecC {
	ignorePadding := func(l *decoder, start Pos) {
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

	// emit should write into output what was read up until this point and move l.start to l.pos
	emit := func(l *decoder, t itemType) {
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

	var (
		startState      stateFn
		invalidState    stateFn
		variantState    stateFn
		alphabetState   stateFn
		urlVariantState stateFn
	)

	startState = func(l *decoder) stateFn {
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

	invalidState = func(l *decoder) stateFn {
		for {
			switch n := l.next(); {
			case bytes.ContainsRune([]byte(alphabet), n):
				l.backup()
				emit(l, itemInvalid)
				return alphabetState

			case bytes.ContainsRune([]byte(variant), n):
				l.backup()
				emit(l, itemInvalid)
				return variantState

			case bytes.ContainsRune([]byte(urlVariant), n):
				l.backup()
				emit(l, itemInvalid)
				return urlVariantState

			case n == -1:
				emit(l, itemInvalid)
				return nil
			}
		}
	}

	//TODO consider ==
	alphabetState = func(l *decoder) stateFn {
		for {
			switch n := l.next(); {
			case bytes.ContainsRune([]byte(alphabet), n):
				l.acceptRun(alphabet)
				continue

			case bytes.ContainsRune([]byte(variant), n):
				l.backup()
				return variantState

			case bytes.ContainsRune([]byte(urlVariant), n):
				l.backup()
				return urlVariantState

			case n == -1:
				emit(l, itemAlphabet)
				return nil

			default:
				l.backup()
				start := l.start
				emit(l, itemAlphabet)
				ignorePadding(l, start)
				return invalidState
			}
		}
	}

	variantState = func(l *decoder) stateFn {
		for {
			switch n := l.next(); {
			case bytes.ContainsRune([]byte(alphabet+variant), n):
				l.acceptRun(alphabet + variant)
				continue

			case n == -1:
				emit(l, itemVariant)
				return nil

			default:
				l.backup()
				start := l.start
				emit(l, itemVariant)
				ignorePadding(l, start)
				return invalidState
			}
		}
	}

	urlVariantState = func(l *decoder) stateFn {
		for {
			switch n := l.next(); {
			case bytes.ContainsRune([]byte(alphabet+urlVariant), n):
				l.acceptRun(alphabet + urlVariant)
				continue

			case n == -1:
				emit(l, itemUrlVariant)
				return nil

			default:
				l.backup()
				start := l.start
				emit(l, itemUrlVariant)
				ignorePadding(l, start)
				return invalidState
			}
		}
	}

	return &Base64{
		dec:   newDecoder(in, startState),
		input: in,
	}
}

func (b *Base64) String() string {
	return b64name
}

func (b *Base64) Decode() (output string) {
	return string(b.dec.decode())
}

func (b *Base64) Encode() (output string) {
	//TODO allow user to decide which encoder
	return base64.StdEncoding.EncodeToString([]byte(b.input))
}

func (b *Base64) Check() (acceptability float64) {
	//TODO redo
	/*
		var c int
		var tot int
		for _, r := range b.input {
			tot++
			if r { //isValid
				c++
			}
		}
		//Heuristic to consider uneven strings as less likely to be valid base64
		if delta := tot % 4; delta != 0 {
			tot += delta
		}
		return float64(c) / float64(tot)
	*/
	return 0
}
