package decode

import (
	"testing"
)

var B64Test = []struct {
	in       string
	eOut     string
	eIsPrint bool
}{
	{
		"Zm9vYmFy",
		"foobar",
		true,
	},
	{
		"Zm9vYm",
		"foob",
		true,
	},
	{
		"Zm9vYm==",
		"foob",
		true,
	},
	{
		"!Zm!Ym",
		genInvalid(1) + "f" + genInvalid(1) + "b",
		false,
	},
	{
		"!Zm9vYm",
		genInvalid(1) + "foob",
		false,
	},

	{
		"Zm9vYmFy.!Zm9vYmFy",
		"foobar" + genInvalid(2) + "foobar",
		false,
	},
	{
		"Zm9vYmF.!Zm9vYmF",
		"fooba" + genInvalid(2) + "fooba",
		false,
	},
	{
		"Zm9vYmF.!Zm9vYmFy.",
		"fooba" + genInvalid(2) + "foobar" + genInvalid(1),
		false,
	},
	{
		"Zm9vYmF.!Zm9vYmF.8",
		"fooba" + genInvalid(2) + "fooba" + genInvalid(2),
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
		"8J+Hq/Cfh7c",
		"\xf0\x9f\x87\xab\xf0\x9f\x87\xb7",
		true,
	},
	{
		"8J-Hq_Cfh7c",
		"\xf0\x9f\x87\xab\xf0\x9f\x87\xb7",
		true,
	},

	{
		"8J+Hq/$8J+Hq/",
		"\xf0\x9f\x87\xab" + genInvalid(1) + "\xf0\x9f\x87\xab",
		false,
	},
	{
		"8J-Hq_Cfh7c$8J-Hq_Cfh7c",
		"\xf0\x9f\x87\xab\xf0\x9f\x87\xb7" + genInvalid(1) + "\xf0\x9f\x87\xab\xf0\x9f\x87\xb7",
		false,
	},
	{
		"$+Hq/Cfh7cg",
		genInvalid(1) + "\xf8z\xbf\t\xf8{r",
		false,
	},
	{
		"$-Hq_Cfh7cg",
		genInvalid(1) + "\xf8z\xbf\t\xf8{r",
		false,
	},
	{
		"+Hq/Cfh7cg",
		"\xf8z\xbf\t\xf8{r",
		false,
	},
	{
		"-Hq_Cfh7cg",
		"\xf8z\xbf\t\xf8{r",
		false,
	},
	{
		"/",
		genInvalid(1),
		false,
	},
	{
		"_",
		genInvalid(1),
		false,
	},
}

var Base64EncodeTest = []struct {
	in   string
	eOut string
}{
	{
		"foobar",
		"Zm9vYmFy",
	},
	{
		"fooba▶︎",
		"Zm9vYmHilrbvuI4=",
	},
	{
		"",
		"",
	},
}

var Base64CheckTest = []struct {
	in   string
	eOut float64
}{
	{
		"Zm9vYmFy",
		1,
	},
	{
		"Zm9vYmF",
		0.7,
	},
}

func TestB64Decode(t *testing.T) {
	for _, tt := range B64Test {
		d := NewB64CodecC(tt.in)
		out := d.Decode()
		if out != tt.eOut {
			t.Errorf("Expected decoded value: '%s' but got '%s'", tt.eOut, out)
		}
		if IsPrint(out) != tt.eIsPrint {
			t.Errorf("Expected printable: %v", tt.eIsPrint)
		}
	}
}

func TestB64Encode(t *testing.T) {
	for _, tt := range Base64EncodeTest {
		d := NewB64CodecC(tt.in)
		out := d.Encode()
		if out != tt.eOut {
			t.Errorf("Expected encoded value: '%s' but got '%s'", tt.eOut, out)
		}
	}
}

func TestB64Check(t *testing.T) {
	for _, tt := range Base64CheckTest {
		d := NewB64CodecC(tt.in)
		out := d.Check()
		if CompareFloat(out, tt.eOut, 0.1) != 0 {
			t.Errorf("Expected check value of '%f' but got '%f'", tt.eOut, out)
		}
	}
}
