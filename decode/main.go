package decode

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/empijei/cli/lg"
)

// MainStandalone parses its own flag and it is the funcion to be run when using
// `wapty decode`. This behaves as a main and expects the "decode" parameter to
// be removed from os.Args.
func MainStandalone(args ...string) {
	buf := takeInput(args)
	sequence := strings.Split(flagCodeclist, ",")
	for _, codec := range sequence {
		//This is to avoid printing twice the final result
		//if i < len(sequence)-2 {
		//fmt.Fprintln(os.Stderr, buf)
		//}
		if out, codecUsed, err := DecodeEncode(buf, flagEncode, codec); err == nil {
			lg.Infof("Codec: %s\n%s\n", codecUsed, out)
			buf = out
		} else {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
	}
}

// DecodeEncode takes an input string `buf` and decodes/encodes it (depending on the
// `encode` parameter) with the given `codec`. It returns the encoded/decoded string
// or an error if the process failed.
func DecodeEncode(buf string, encode bool, codec string) (out string, codecUsed string, err error) {

	// Build list of available codecs
	var codecNames []string
	for _, cc := range codecs {
		codecNames = append(codecNames, cc.name)
	}
	codecNamesStr := strings.Join(codecNames, ", ")

	var c CodecC
	if codec == "smart" {
		if encode {
			err = fmt.Errorf("Cannot 'smart' encode, please specify a codec")
			return
		}
		c = SmartDecode(buf)
	} else {
		for _, cc := range codecs {
			if cc.name == codec {
				c = cc.codecCons(buf)
			}
		}
		if c == nil {
			err = fmt.Errorf("Codec not found: '%s'. Supported codecs are: %s\n", codec, codecNamesStr)
			return
		}
	}
	codecUsed = c.Name()
	if encode {
		out = c.Encode()
	} else {
		out = c.Decode()
	}
	return
}

func takeInput(args []string) string {
	stdininfo, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while connecting to stdin: %s\n", err.Error())
	}
	if err == nil && stdininfo.Mode()&os.ModeCharDevice == 0 {
		//The input is a pipe, so I assume it is what I'm going to decode/encode
		fmt.Fprintln(os.Stderr, "Reading from stdin...")
		buf, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while reading from stdin: %s\n", err.Error())
			os.Exit(2)
		}
		return string(buf)
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Didn't find anything to decode/encode, exiting...")
		os.Exit(2)
	}
	return args[0]
}
