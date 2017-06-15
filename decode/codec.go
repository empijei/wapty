package decode

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"unicode"
)

type codecConstructor func(string) CodecC

var codecs = make(map[string]codecConstructor)
var codecsM sync.Mutex

func addCodecC(name string, c codecConstructor) {
	codecsM.Lock()
	defer codecsM.Unlock()
	codecs[name] = c
}

type Decoder interface {
	//Decodes the string and returns a decoded value that tries to skip invalid
	//input and to decode as much as possible.
	//Returns if the decoded string can be printed as valid unicode.
	Decode() (output string)
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

// IsPrint checks if a decoded string is a valid utf string
func IsPrint(decoded string) bool {
	if strings.Contains(decoded, string(invalid)) {
		return false
	}
	for _, r := range decoded {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
