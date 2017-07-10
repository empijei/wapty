package decode

import (
	"strings"
	"testing"
)

var Base16Test = []struct {
	in       string
	eOut     string
	eIsPrint bool
}{
	{
		"666F6F626172",
		"foobar",
		true,
	},
	{
		"666f6f626172",
		"foobar",
		true,
	},
	{
		"666F6F62617",
		"fooba" + genInvalid(1),
		false,
	},
	{
		"666F6F626172.!666F6F626172",
		"foobar" + genInvalid(2) + "foobar",
		false,
	},
	{
		"666F6F62617.!666F6F62617",
		"fooba" + genInvalid(3) + "fooba" + genInvalid(1),
		false,
	},
	{
		"666F6F62617.!666F6F62617.",
		"fooba" + genInvalid(3) + "fooba" + genInvalid(2),
		false,
	},
	{
		"666F6F62617.!666F6F62617.8",
		"fooba" + genInvalid(3) + "fooba" + genInvalid(3),
		false,
	},
	{
		"6",
		genInvalid(1),
		false,
	},
	{
		"",
		"",
		true,
	},
	{
		".",
		genInvalid(1),
		false,
	},
}

var Base16EncodeTest = []struct {
	in   string
	eOut string
}{
	{
		"foobar",
		"666F6F626172",
	},
	{
		"fooba▶︎",
		"666F6F6261E296B6EFB88E",
	},
	{
		"",
		"",
	},
}

var Base16CheckTest = []struct {
	in   string
	eOut float64
}{
	{
		"666F6F626172",
		1,
	},
	{
		"666F6F62617",
		0.91,
	},
}

func TestB16Decode(t *testing.T) {
	for _, tt := range Base16Test {
		d := NewB16CodecC(tt.in)
		out := d.Decode()
		if out != tt.eOut {
			t.Errorf("Expected decoded value: '%s' but got '%s'", tt.eOut, out)
		}
		if IsPrint(out) != tt.eIsPrint {
			t.Errorf("Expected printable: %v", tt.eIsPrint)
		}
	}
}

func TestB16Encode(t *testing.T) {
	for _, tt := range Base16EncodeTest {
		d := NewB16CodecC(tt.in)
		out := d.Encode()
		if strings.ToUpper(out) != tt.eOut {
			t.Errorf("Expected encoded value: '%s' but got '%s'", tt.eOut, out)
		}
	}
}

func TestB16Check(t *testing.T) {
	for _, tt := range Base16CheckTest {
		d := NewB16CodecC(tt.in)
		out := d.Check()
		if CompareFloat(out, tt.eOut, 0.1) != 0 {
			t.Errorf("Expected check value of '%f' but got '%f'", tt.eOut, out)
		}
	}
}

func CompareFloat(a, b float64, tolerance float64) int {
	if a < b-tolerance {
		return 1
	}
	if b > a+tolerance {
		return -1
	}
	return 0
}
