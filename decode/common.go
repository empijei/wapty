package decode

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"unicode"
)

var invalid = '�'

type codecConstructor func(string) CodecC

var codecs = make(map[string]codecConstructor)
var codecs_m sync.Mutex

func addCodecC(name string, c codecConstructor) {
	codecs_m.Lock()
	defer codecs_m.Unlock()
	codecs[name] = c
}

func genInvalid(n int) (inv string) {
	return strings.Repeat(string(invalid), n)
}

type Decoder interface {
	//Decodes the string and returns a decoded value that tries to skip invalid
	//input and to decode as much as possible.
	//Returns if the decoded string can be printed as valid unicode.
	Decode() (output string, isPrintable bool)
}

type Encoder interface {
	//Encodes the string
	Encode() (output string)
}

type Checker interface {
	//Returns a metric to determine how likely it is for the given string to be
	//a valid value for the specified Checker Type.
	//The likelihood always ranges between 0 and 1
	Check() (acceptability float64)
}

type CodecC interface {
	Decoder
	Encoder
	Checker
	fmt.Stringer
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

func SmartDecode(input string) (c CodecC) {
	//loop through the available CodecCs and determine which one is the best one
	var curvalue float64
	for _, cc := range codecs {
		tmp := cc(input)
		if t := tmp.Check(); t > curvalue {
			curvalue = t
			c = tmp
		}
	}
	log.Printf("Smart Decoding, selected: %s with likelihood==%d%%", c.String(), int(curvalue*100))
	return
}
