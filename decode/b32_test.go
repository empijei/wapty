package decode

import (
	"strings"
	"testing"
)

var base32Test = []struct {
	in       string
	eOut     string
	eIsPrint bool
}{
	{
		"GFTG633CMFZA====",
		"1foobar",
		true,
	},
	{
		"mzxw6ytboi======",
		"foobar",
		true,
	},
	{
		"MZXW6YTBO",
		"fooba" + genInvalid(1),
		false,
	},
	{
		"MZXW6YTBOI.!MZXW6YTBOI",
		"foobar" + genInvalid(2) + "foobar",
		false,
	},
	{
		"MZXW6YTBO.!MZXW6YTBO",
		"fooba" + genInvalid(3) + "fooba" + genInvalid(1),
		false,
	},
	{
		"MZXW6YTBO.!MZXW6YTBO.",
		"fooba" + genInvalid(3) + "fooba" + genInvalid(2),
		false,
	},
	{
		"MZXW6YTBO.!MZXW6YTBO.8",
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
}

var base32EncodeTest = []struct {
	in   string
	eOut string
}{
	{
		"foobar",
		"MZXW6YTBOI======",
	},
	{
		"fooba▶︎",
		"MZXW6YTB4KLLN35YRY======",
	},
	{
		"",
		"",
	},
}

var base32CheckTest = []struct {
	in   string
	eOut float64
}{
	{
		"MZXW6YTBOI",
		1,
	},
	{
		"MZXW6YTBO",
		0.91,
	},
}

func TestB32Decode(t *testing.T) {
	for _, tt := range base32Test {
		d := NewB32CodecC(tt.in)
		out := d.Decode()
		if out != tt.eOut {
			t.Errorf("Expected decoded value: '%s' but got '%s'", tt.eOut, out)
		}
		if IsPrint(out) != tt.eIsPrint {
			t.Errorf("Expected printable: %v", tt.eIsPrint)
		}
	}
}

func TestB32Encode(t *testing.T) {
	for _, tt := range base32EncodeTest {
		d := NewB32CodecC(tt.in)
		out := d.Encode()
		if strings.ToUpper(out) != tt.eOut {
			t.Errorf("Expected encoded value: '%s' but got '%s'", tt.eOut, out)
		}
	}
}

func TestB32Check(t *testing.T) {
	for _, tt := range base32CheckTest {
		d := NewB32CodecC(tt.in)
		out := d.Check()
		if CompareFloat(out, tt.eOut, 0.1) != 0 {
			t.Errorf("Expected check value of '%f' but got '%f'", tt.eOut, out)
		}
	}
}
