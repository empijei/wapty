package decode

import "testing"

var URLDecodeTest = []struct {
	in   string
	eOut string
}{
	{
		"foo%20bar",
		"foo bar",
	},
	{
		"%C3%BCnter",
		"ünter",
	},
	{
		"",
		"",
	},
}

var URLEncodeTest = []struct {
	in   string
	eOut string
}{
	{
		"foo bar",
		"foo%20bar",
	},
	{
		"ünter",
		"%C3%BCnter",
	},
	{
		"",
		"",
	},
}

var URLCheckTest = []struct {
	in   string
	eOut float64
}{
	{
		"foo%20ba",
		1,
	},
	{
		"foo%2Xba",
		0.875,
	},
	{
		"foo%2",
		0.8,
	},
	{
		"",
		0.0,
	},
}

func TestURLDecode(t *testing.T) {
	for _, tt := range URLDecodeTest {
		d := NewURLCodecC(tt.in)
		out := d.Decode()
		if out != tt.eOut {
			t.Errorf("Expected decoded value: '%s' but got '%s'", tt.eOut, out)
		}
	}
}

func TestURLEncode(t *testing.T) {
	for _, tt := range URLEncodeTest {
		d := NewURLCodecC(tt.in)
		out := d.Encode()
		if out != tt.eOut {
			t.Errorf("Expected encoded value: '%s' but got '%s'", tt.eOut, out)
		}
	}
}

func TestURLCheck(t *testing.T) {
	for _, tt := range Base16CheckTest {
		d := NewURLCodecC(tt.in)
		out := d.Check()
		if CompareFloat(out, tt.eOut, 0.1) != 0 {
			t.Errorf("Expected check value of '%f' but got '%f'", tt.eOut, out)
		}
	}
}
