package decode

import (
	"bytes"
	"net/url"
)

const urlname = "url"

// URL takes an input string
type URL struct {
	input string
}

// NewURLCodecC state machine to smartly decode a string with invalid chars
func NewURLCodecC(in string) CodecC {
	return &URL{
		input: in,
	}
}

// Name returns the name of the codec
func (b *URL) Name() string {
	return urlname
}

// Decode a valid url encoded string
func (b *URL) Decode() (output string) {
	res, err := url.PathUnescape(b.input)
	if err != nil {
		return "Not a valid URL encoded string"
	}
	return res
}

// Encode a string to url encode
func (b *URL) Encode() (output string) {
	return url.PathEscape(b.input)
}

// Check returns the percentage of valid url characters in the input string
func (b *URL) Check() (acceptability float64) {
	var c int
	var tot int
	for pos, char := range b.input {
		tot++
		if bytes.ContainsRune([]byte("%"), char) && pos < len(b.input)+2 {
			if bytes.ContainsAny([]byte(b16Alphabet), string(b.input[pos+1])) &&
				bytes.ContainsAny([]byte(b16Alphabet), string(b.input[pos+2])) {
				c++
			}
		}
		if bytes.ContainsRune([]byte(b64Alphabet), char) {
			c++
		}
	}
	return float64(c) / float64(tot)
}
