package decode

import (
	"fmt"
	"strings"
	"testing"
)

type expectation int

const (
	expectGood expectation = iota // expect a good result with no errors
	expectBad                     // expect a bad result with no errors
	expectErr                     // expect errors
)

var testDecodeData = []struct {
	cFlag  string
	in     string
	eOut   string
	expect expectation
}{
	{
		"b16",
		"666F6F626172",
		"foobar",
		expectGood,
	},
	{
		"b16,b64",
		"5a6d3976",
		"foo",
		expectGood,
	},
	{
		",,",
		"5a6d3976",
		"",
		expectErr,
	},
	{
		"",
		"5a6d3976",
		"",
		expectErr,
	},
	{
		"smart,smart",
		"5a6d3976",
		"foo",
		expectGood,
	},
}

func TestDecode(t *testing.T) {
	succFailStr := map[bool]string{
		true:  "good result",
		false: "bad result",
	}
	for _, test := range testDecodeData {
		wasErrRaised := false
		buf := test.in
		for _, codec := range strings.Split(test.cFlag, ",") {
			var err error
			buf, _, err = DecodeEncode(buf, false, codec)
			if err != nil {
				if test.expect != expectErr {
					t.Errorf("On call DecodeEncode('%s', false, '%s'):\n", buf, codec)
					t.Fatalf("Unexpected error while decoding: %s", err.Error())
				}
				wasErrRaised = true
			}
		}
		if test.expect == expectErr {
			if !wasErrRaised {
				t.Errorf("Expected error on test decode(in='%s', codecs='%s') -> '%s', but got none",
					test.in, test.cFlag, test.eOut)
			}
		} else if (buf == test.eOut) != (test.expect == expectGood) {
			t.Errorf("Expected %s on test decode(in='%s', codecs='%s') -> '%s', but got %s",
				succFailStr[test.expect == expectGood],
				test.in,
				test.cFlag,
				test.eOut,
				succFailStr[buf == test.eOut],
			)
			t.Errorf("\n\tout: '%s'\n\n\texpected: '%s'\n", buf, test.eOut)
		}
	}
}

var testGoodCodecData = []struct {
	in         string
	codecsUsed []string
}{
	{
		"Zm9vYmFyIC1uCg==",
		[]string{"b64"},
	},
	{
		"5a6d3976",
		[]string{"b16", "b64"},
	},
}

func TestGoodCodec(t *testing.T) {
	for _, test := range testGoodCodecData {
		buf := test.in
		var codecUsed string
		var err error
		for _, expectedCodec := range test.codecsUsed {
			prevBuf := buf
			buf, codecUsed, err = DecodeEncode(buf, false, "smart")
			if err != nil {
				t.Errorf("On call DecodeEncode('%s', false, 'smart'):\n", buf)
				t.Fatalf("Unexpected error while decoding: %s", err.Error())
			}
			if codecUsed != expectedCodec {
				t.Errorf("Decoding '%s': expected codec '%s' but got '%s'",
					prevBuf, expectedCodec, codecUsed)
			}
		}
	}
}

func TestDecodeStdin(t *testing.T) {
	// TODO
}

var testEncodeData = []struct {
	in         string
	codecsUsed []string
	out        string
}{
	{
		"foobar",
		[]string{"b64"},
		"Zm9vYmFy",
	},
	{
		"asdasd",
		[]string{"b16"},
		"617364617364",
	},
	{
		"wapty",
		[]string{"b32"},
		"O5QXA5DZ",
	},
	{
		"asdasd",
		[]string{"b16", "b32", "b64"},
		"R1lZVE9NWldHUTNEQ05aVEdZMkE9PT09",
	},
	{
		"asdasd",
		[]string{"b16", "b16", "b16"},
		"333633313337333333363334333633313337333333363334",
	},
	{
		"asdasd",
		[]string{"b32", "b32", "b32"},
		"JJLEIRSVKYZEUTCGI5DESVCLKJEFKNSUGJIEUNKIKU6T2PJ5HU6Q====",
	},
	{
		"asdasd",
		[]string{"b64", "b64", "b64"},
		"V1ZoT2ExbFlUbXM9",
	},
}

func TestEncode(t *testing.T) {
	type step struct {
		codec string
		out   string
	}
	for _, test := range testEncodeData {
		buf := test.in
		var err error
		var intermediate []step
		for _, codecUsed := range test.codecsUsed {
			buf, _, err = DecodeEncode(buf, true, codecUsed)
			intermediate = append(intermediate, step{codecUsed, buf})
			if err != nil {
				t.Errorf("On call DecodeEncode('%s', true, '%s'):\n", buf, codecUsed)
				t.Fatalf("Unexpected error while encoding: %s", err.Error())
			}
		}
		if buf != test.out {
			t.Errorf("Encoding '%s' with %s:\n\n\texpected\n\t\t'%s'\n\tbut got\n\t\t'%s'",
				test.in, strings.Join(test.codecsUsed, ","), test.out, buf)
			s := "Intermediate steps:\n"
			for _, stp := range intermediate {
				s += fmt.Sprintf("\tcodec: %s\n\tresult: %s\n", stp.codec, stp.out)
			}
			t.Errorf(s)
		}
	}
}

func TestEncodeStdin(t *testing.T) {
	// TODO
}
