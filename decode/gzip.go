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
		fmt.Println("Secondo errore")
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
	// TODO
	return
}
