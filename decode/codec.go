package decode

import (
	"log"
	"strings"
	"unicode"
)

type codecConstructor func(string) CodecC

var codecs = []struct {
	name      string
	codecCons codecConstructor
}{
	{
		b16name,
		codecConstructor(NewB16CodecC),
	},
	{
		b32name,
		codecConstructor(NewB32CodecC),
	},
	{
		b64name,
		codecConstructor(NewB64CodecC),
	},
}

//Decoder decodes the string and returns a decoded value that tries to skip
//invalid input and to decode as much as possible.
//Returns if the decoded string can be printed as valid unicode.
type Decoder interface {
	Decode() (output string)
}

//Encoder encodes the string
type Encoder interface {
	Encode() (output string)
}

//Checker returns a metric to determine how likely it is for the given string
//to be a valid value for the specified Checker Type.
//The likelihood always ranges between 0 and 1
type Checker interface {
	Check() (acceptability float64)
}

//CodecC creates an interface of interfaces usable by other codecs
type CodecC interface {
	Decoder
	Encoder
	Checker
	Name() string
}

//SmartDecode loops through the available CodecCs
//and determine which one is the best one to use
func SmartDecode(input string) (c CodecC) {
	var curvalue float64
	//FIXME add a null codecC if no codecC is selected
	for _, cc := range codecs {
		tmp := cc.codecCons(input)
		if t := tmp.Check(); t > curvalue {
			curvalue = t
			c = tmp
		}
	}
	log.Printf("Smart Decoding, selected: %s with likelihood==%d%%", c.Name(), int(curvalue*100))
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
