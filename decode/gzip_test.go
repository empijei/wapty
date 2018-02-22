package decode

import (
	"fmt"
	"testing"
)

var GzipDecodeTest = []struct {
	in   string
	eOut string
}{
	{
		"H4sICAxCLFoAA3RvY29tcHJlc3MAc8vPd0os4gIA5QGj3QcAAAA=",
		"Name: tocompress\nFooBar\n",
	},
	{
		"",
		"",
	},
}

var GzipCheckTest = []struct {
	in   string
	eOut float64
}{
	{
		// 1f8b08080c422c5a0003746f636f6d70726573730073cbcf774a2ce20200e501a3dd07000000139
		"H4sICAxCLFoAA3RvY29tcHJlc3MAc8vPd0os4gIA5QGj3QcAAAA=",
		1,
	},
	{
		"H4sI",
		0,
	},
}

func TestGzipDecode(t *testing.T) {
	for _, tt := range GzipDecodeTest {
		b64 := NewB64CodecC(tt.in)
		input := b64.Decode()
		d := NewGzipCodecC(input)
		out := d.Decode()
		if out != tt.eOut {
			fmt.Println(out)
			t.Errorf("Expected decoded value: '%s' but got '%s'", tt.eOut, out)
		}
	}
}

func TestGzipCheck(t *testing.T) {
	for _, tt := range GzipCheckTest {
		b64 := NewB64CodecC(tt.in)
		input := b64.Decode()
		d := NewGzipCodecC(input)
		out := d.Check()
		if CompareFloat(out, tt.eOut, 0.1) != 0 {
			t.Errorf("Expected acceptability value: %f, but got %f", tt.eOut, out)
		}
	}
}
