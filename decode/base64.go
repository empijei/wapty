package decode

import (
	"bytes"
	"encoding/base64"
	"strings"
)

type Base64 struct {
	input  string
	cursor int
	pos    int
	output *bytes.Buffer
}

//TODO recognise other encodings
const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const variant = "+/"
const urlVariant = "-_"

var encoding = base64.StdEncoding.WithPadding(base64.NoPadding)

func NewBase64CodecC(in string) *Base64 {
	return &Base64{
		input:  in,
		output: bytes.NewBuffer(make([]byte, 0, encoding.DecodedLen(len(in)))),
	}
}

//TODO modificare in modo che prenda pure senza padding
func (b *Base64) nextValid() {
	validseen := 0
	for b.pos < len(b.input) &&
		validseen < 4 {
		if b.isValid(rune(b.input[b.pos])) {
			validseen++
			if validseen == 4 {
				b.pos -= 4
			}
		} else {
			validseen = 0
		}
		b.pos++
	}
}

func (b *Base64) acceptRun() {
	for b.pos < len(b.input) && b.isValid(rune(b.input[b.pos])) {
		b.pos++
	}
	if delta := (b.pos - b.cursor) % 4; delta != 0 && b.pos < len(b.input) {
		b.pos -= delta
	}
}

func (b *Base64) decodeChunk() {
	buf, err := encoding.DecodeString(b.input[b.cursor:b.pos])
	if err != nil {
		panic("Error when less expected: " + err.Error())
	}
	_, _ = b.output.Write(buf)
	b.cursor = b.pos
}

func (b *Base64) isValid(r rune) bool {
	return strings.ContainsAny(string(r), alphabet+variant)
}

func (b *Base64) Decode() (output string, isPrintable bool) {
	out, err := encoding.DecodeString(b.input)
	if err != nil {
		//Decode as much as possible
		for b.cursor < len(b.input) {
			b.acceptRun()
			b.decodeChunk()
			b.nextValid()
			b.output.WriteString(genInvalid(b.pos - b.cursor))
			b.cursor = b.pos
		}
		output = string(b.output.Bytes())
	} else {
		output = string(out)
	}
	isPrintable = isStringPrintable(output)
	return
}

func (b *Base64) Encode() (output string) {
	return encoding.EncodeToString([]byte(b.input))
}

func (b *Base64) Check() (acceptability float64) {
	//TODO use cursor
	var c int
	var tot int
	for _, r := range b.input {
		tot++
		if b.isValid(r) {
			c++
		}
	}
	//Heuristic to consider uneven strings as less likely to be valid base64
	if delta := tot % 4; delta != 0 {
		tot += delta
	}
	return float64(c) / float64(tot)
}
