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
