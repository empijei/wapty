package decode

import (
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
