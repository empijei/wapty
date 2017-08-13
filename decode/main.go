package decode

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// MainStandalone parses its own flag and it is the funcion to be run when using
// `wapty decode`. This behaves as a main and expects the "decode" parameter to
// be removed from os.Args.
func MainStandalone(args ...string) {

	// FIXME: argument validation should be separated from encoding/decoding
	buf := takeInput(args)
	sequence := strings.Split(flagCodeclist, ",")
	for i, codec := range sequence {
		var c CodecC
		var codecNames []string
		if codec == "smart" {
			if flagEncode {
				fmt.Fprintf(os.Stderr, "Cannot 'smart' encode, please specify a codec")
				os.Exit(2)
			}
			c = SmartDecode(buf)
		} else {
			for _, cc := range codecs {
				if cc.name == codec {
					c = cc.codecCons(buf)
				}
				codecNames = append(codecNames, cc.name)
			}
			if c == nil {
				fmt.Fprintf(os.Stderr, "Codec not found: %s. Supported codecs are: %s\n", codec, strings.Join(codecNames, ", "))
				os.Exit(2)
			}
		}
		fmt.Fprintf(os.Stderr, "Codec: %s\n", c.Name())
		if flagEncode {
			buf = c.Encode()
		} else {
			buf = c.Decode()
		}
		//This is to avoid printing twice the final result
		if i < len(sequence)-2 {
			fmt.Fprintln(os.Stderr, buf)
		}
	}
	fmt.Printf(buf)
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
