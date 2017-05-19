package decode

import (
	"bytes"
	"encoding/hex"
	"strings"
)

type Base16 struct {
	input  string
	cursor int
	pos    int
	output *bytes.Buffer
}

func NewBase16CodecC(in string) *Base16 {
	return &Base16{
		input:  in,
		output: bytes.NewBuffer(make([]byte, 0, hex.DecodedLen(len(in)))),
	}
}

func (b *Base16) nextValid() {
	//FIXME ignore a single character followed by EOF or invalid
	validseen := 0
	for b.pos < len(b.input) &&
		validseen < 2 {
		if b.isValid(rune(b.input[b.pos])) {
			validseen++
			if validseen == 2 {
				b.pos -= 2
			}
		} else {
			validseen = 0
		}
		b.pos++
	}
}

func (b *Base16) acceptRun() {
	for b.pos < len(b.input) && b.isValid(rune(b.input[b.pos])) {
		b.pos++
	}
	if (b.pos-b.cursor)%2 != 0 {
		b.pos--
	}
	//TODO: backup if odd pos-cur
}

func (b *Base16) decodeChunk() {
	buf, err := hex.DecodeString(b.input[b.cursor:b.pos])
	if err != nil {
		panic("Error when less expected: " + err.Error())
	}
	_, _ = b.output.Write(buf)
	b.cursor = b.pos
}

func (b *Base16) isValid(r rune) bool {
	return strings.ContainsAny(string(r), "0123456789abcdefABCDEF")
}

func (b *Base16) Decode() (output string, isPrintable bool) {
	out, err := hex.DecodeString(b.input)
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

func (b *Base16) Encode() (output string) {
	return hex.EncodeToString([]byte(b.input))
}

func (b *Base16) Check() (acceptability float64) {
	//TODO use cursor
	var c int
	var tot int
	for _, r := range b.input {
		tot++
		if b.isValid(r) {
			c++
		}
	}
	//Heuristic to consider uneven strings as less likely to be velid hex
	if tot%2 == 0 {
		tot++
	}
	return float64(c) / float64(tot)
}
