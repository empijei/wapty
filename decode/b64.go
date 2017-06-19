package decode

import (
	"bytes"
	"encoding/base64"
)

const b64Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const b64Variant = "+/"
const b64Padding = "="
const b64UrlVariant = "-_"

const b64name = "b64"

func init() {
	addCodecC(b64name, codecConstructor(NewB64CodecC))
}

type Base64 struct {
	dec   *decoder
	input string
}

// nolint: gocyclo
func NewB64CodecC(in string) CodecC {
	const (
		itemInvalid itemType = iota
		itemAlphabet
		itemVariant
		itemUrlVariant
	)

	ignorePadding := func(d *decoder, start Pos) {
		for {
			if d.peek() != '=' {
				return
			}
			switch (d.pos - start) % 4 {
			case 2, 3:
				d.next()
				d.ignore()

			default:
				return
			}
		}
	}

	// emit should write into output what was read up until this point and move l.start to l.pos
	emit := func(d *decoder, t itemType) {
		token := d.input[d.start:d.pos]

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

		d.out.Write(decodefunc(token))
		d.start = d.pos
	}

	var (
		startState      stateFn
		invalidState    stateFn
		variantState    stateFn
		alphabetState   stateFn
		urlVariantState stateFn
	)

	startState = func(d *decoder) stateFn {
		switch n := d.peek(); {
		case bytes.ContainsRune([]byte(b64Alphabet), n):
			return alphabetState
		case bytes.ContainsRune([]byte(b64Variant), n):
			return variantState
		case bytes.ContainsRune([]byte(b64UrlVariant), n):
			return urlVariantState
		case n == eof:
			return nil
		default:
			return invalidState
		}
	}

	invalidState = func(d *decoder) stateFn {
		for {
			switch n := d.next(); {
			case bytes.ContainsRune([]byte(b64Alphabet), n):
				d.backup()
				emit(d, itemInvalid)
				return alphabetState

			case bytes.ContainsRune([]byte(b64Variant), n):
				d.backup()
				emit(d, itemInvalid)
				return variantState

			case bytes.ContainsRune([]byte(b64UrlVariant), n):
				d.backup()
				emit(d, itemInvalid)
				return urlVariantState

			case n == eof:
				emit(d, itemInvalid)
				return nil
			}
		}
	}

	alphabetState = func(d *decoder) stateFn {
		for {
			switch n := d.next(); {
			case bytes.ContainsRune([]byte(b64Alphabet), n):
				d.acceptRun(b64Alphabet)
				continue

			case bytes.ContainsRune([]byte(b64Variant), n):
				d.backup()
				return variantState

			case bytes.ContainsRune([]byte(b64UrlVariant), n):
				d.backup()
				return urlVariantState

			case n == eof:
				emit(d, itemAlphabet)
				return nil

			default:
				d.backup()
				start := d.start
				emit(d, itemAlphabet)
				ignorePadding(d, start)
				return invalidState
			}
		}
	}

	variantState = func(d *decoder) stateFn {
		for {
			switch n := d.next(); {
			case bytes.ContainsRune([]byte(b64Alphabet+b64Variant), n):
				d.acceptRun(b64Alphabet + b64Variant)
				continue

			case n == eof:
				emit(d, itemVariant)
				return nil

			default:
				d.backup()
				start := d.start
				emit(d, itemVariant)
				ignorePadding(d, start)
				return invalidState
			}
		}
	}

	urlVariantState = func(d *decoder) stateFn {
		for {
			switch n := d.next(); {
			case bytes.ContainsRune([]byte(b64Alphabet+b64UrlVariant), n):
				d.acceptRun(b64Alphabet + b64UrlVariant)
				continue

			case n == eof:
				emit(d, itemUrlVariant)
				return nil

			default:
				d.backup()
				start := d.start
				emit(d, itemUrlVariant)
				ignorePadding(d, start)
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
	var c int
	var tot int
	for _, r := range b.input {
		tot++
		if bytes.ContainsRune([]byte(b64Alphabet+b64Variant+b64UrlVariant+b64Padding), r) {
			c++
		}
	}
	//Heuristic to consider uneven strings as less likely to be valid base64
	if delta := tot % 4; delta != 0 {
		tot += delta
	}
	return float64(c) / float64(tot)
}
