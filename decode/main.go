package decode

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func MainStandalone() {
	encode := flag.Bool("encode", false, "Sets the decoder to an encoder instead")
	codeclist := flag.String("codec", "smart", "Sets the decoder/encoder codec. Multiple codecs can be specified and comma separated, they will be applied one on the output of the previous as in a pipeline.")
	flag.Parse()

	buf := takeInput()
	sequence := strings.Split(*codeclist, ",")
	for i, codec := range sequence {
		var c CodecC
		if codec == "smart" {
			if *encode {
				fmt.Fprintf(os.Stderr, "Cannot 'smart' encode, please specify a codec")
				os.Exit(2)
			}
			c = SmartDecode(buf)
		} else {
			var names []string
			for name, cc := range codecs {
				if name == codec {
					c = cc(buf)
				}
			}
			if c == nil {
				fmt.Fprintf(os.Stderr, "Codec not found: %s. Supported codec are: %s", codec, strings.Join(names, ", "))
				os.Exit(2)
			}
		}
		fmt.Fprintf(os.Stderr, "Codec: %s\n", c.String())
		if *encode {
			buf = c.Encode()
		} else {
			buf, _ = c.Decode()
		}
		//This is to avoid printing twice the final result
		if i < len(sequence)-2 {
			fmt.Fprintln(os.Stderr, buf)
		}
	}
	fmt.Printf(buf)
}

func takeInput() string {
	args := flag.Args()
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
	} else {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Didn't find anything to decode/encode, exiting...")
			os.Exit(2)
		}
		return args[0]
	}
}
