package decode

import (
	"testing"
)

var Base64Test = []struct {
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
		"Zm9vYm==",
		"foob",
		true,
	},
	{
		"Zm9vYmFy.!Zm9vYmFy",
		"foobar" + genInvalid(2) + "foobar",
		false,
	},
	//{
	//"Zm9vYmF.!Zm9vYmF",
	//"fooba" + genInvalid(3) + "fooba" + genInvalid(1),
	//false,
	//},
	//{
	//"Zm9vYmF.!Zm9vYmFy.",
	//"fooba" + genInvalid(3) + "fooba" + genInvalid(2),
	//false,
	//},
	//{
	//"Zm9vYmF.!Zm9vYmF.8",
	//"fooba" + genInvalid(3) + "fooba" + genInvalid(3),
	//false,
	//},
	//{
	//"6",
	//genInvalid(1),
	//false,
	//},
}

func TestB64Decode(t *testing.T) {
	//invalid = '§'
	for _, tt := range Base64Test {
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
