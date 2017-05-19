package decode

import (
	"strings"
	"unicode"
)

var invalid = '�'

func genInvalid(n int) (inv string) {
	return strings.Repeat(string(invalid), n)
}

type Decoder interface {
	//Decodes the string and returns a decoded value that tries to skip invalid
	//input and to decode as much as possible.
	//Returns if the decoded string can be printed as valid unicode.
	Decode(input string) (output string, isPrintable bool)
}

type Encoder interface {
	//Encodes the string
	Encode(input string) (output string)
}

type Checker interface {
	//Returns a metric to determine how likely it is for the given string to be
	//a valid value for the specified Checker Type.
	//The likelyhood always ranges between 0 and 1
	Check(input string) (acceptability float64)
}

type CodecC interface {
	Decoder
	Encoder
	Checker
}

func isStringPrintable(s string) bool {
	if strings.Contains(s, string(invalid)) {
		return false
	}
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func SmartDecode(input string) (output string, isPrintable bool) {
	//loop through the available CodecCs and determine which one is the best one
	panic("Not implemented yet")
}
