package decode

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"strings"
)

const gzipname = "gzip"

// gzip takes an input string
type gzipDec struct {
	input string
}

func NewGzipCodecC(in string) CodecC {
	return &gzipDec{
		input: in,
	}
}

// Name returns the name of the codec
func (b *gzipDec) Name() string {
	return gzipname
}

// Decode a valid gzip compressed string
func (b *gzipDec) Decode() (output string) {
	buf := new(bytes.Buffer)

	zr, err := gzip.NewReader(strings.NewReader(b.input))

	if err == io.EOF {
		return ""
	}

	if err != nil {
		log.Fatal(err)
	}

	buf.WriteString(fmt.Sprintf("Name: %s\n", zr.Name))

	if com := zr.Comment; com != "" {
		buf.WriteString(fmt.Sprintf("Comment: %s\n", com))
	}

	if _, err := io.Copy(buf, zr); err != nil {
		return "Not a valid gzip compressed string"
	}

	if err := zr.Close(); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

// Encode compresses a string with gzip
func (b *gzipDec) Encode() (output string) {
	return //TODO
}

// Check returns the probability a string is gzip compressed
func (b *gzipDec) Check() (acceptability float64) {
	var c int
	tot := 8
	gzipID1 := 0x1f
	gzipID2 := 0x8b
	buf := []byte(b.input)

	if len(buf) < 10 {
		return 0.
	}

	// this is the first ID byte
	if buf[0] == byte(gzipID1) {
		c = c + 2
	}

	// this is the secondo ID byte
	if buf[1] == byte(gzipID2) {
		c = c + 2
	}

	// this is the compression method
	if buf[2] <= 8 {
		c = c + 2
	}

	// this is the flags
	if buf[3] < 32 {
		c++
	}

	// this is the operating system
	if buf[9] <= 13 || buf[9] == 255 {
		c++
	}
	return float64(c) / float64(tot)
}
