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

func TestBase64Decode(t *testing.T) {
	//FIXME empijei: this test was disabled
	return
	//invalid = '§'
	for _, tt := range B64Test {
		d := NewBase64CodecC(tt.in)
		out, ip := d.Decode()
		if out != tt.eOut {
			t.Errorf("Expected decoded value: '%s' but got '%s'", tt.eOut, out)
		}
		if ip != tt.eIsPrint {
			t.Errorf("Expected printable: %v but got %v", tt.eIsPrint, ip)
		}
	}
}
