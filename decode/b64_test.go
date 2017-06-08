package decode

import (
	"bytes"
	"testing"
)

var Base64Test = []struct {
	in   string
	eOut string
}{
	{
		"Zm9vYmFy",
		"foobar",
	},
	{
		"Zm9vYm",
		"foob",
	},
	{
		"Zm9vYm==",
		"foob",
	},
	{
		"!Zm!Ym",
		genInvalid(1) + "f" + genInvalid(1) + "b",
	},
	{
		"!Zm9vYm",
		genInvalid(1) + "foob",
	},

	{
		"Zm9vYmFy.!Zm9vYmFy",
		"foobar" + genInvalid(2) + "foobar",
	},
	{
		"Zm9vYmF.!Zm9vYmF",
		"fooba" + genInvalid(2) + "fooba",
	},
	{
		"Zm9vYmF.!Zm9vYmFy.",
		"fooba" + genInvalid(2) + "foobar" + genInvalid(1),
	},
	{
		"Zm9vYmF.!Zm9vYmF.8",
		"fooba" + genInvalid(2) + "fooba" + genInvalid(2),
	},
	{
		"6",
		genInvalid(1),
	},
}

func TestB64Decode(t *testing.T) {
	//invalid = '§'
	for _, tt := range Base64Test {
		d := NewB64Decoder(tt.in)
		out := d.decode()
		if bytes.Compare(out, []byte(tt.eOut)) != 0 {
			t.Errorf("Expected decoded value: '%s' but got '%s'", tt.eOut, out)
		}
	}
}
