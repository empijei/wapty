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
		"H4sICAxCLFoAA3RvY29tcHJlc3MAc8vPd0os4gIA5QGj3QcAAAA=",
		1,
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
		d := NewGzipCodecC(tt.in)
		out := d.Check()
		if CompareFloat(out, tt.eOut, 0.1) != 0 {
			t.Errorf("Expected acceptability value: %f, but got %f", tt.eOut, out)
		}
	}
}
