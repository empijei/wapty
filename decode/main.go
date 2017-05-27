package decode

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func MainStandalone(codec string, encode bool) {
	buf := takeInput()
	var c CodecC
	if codec == "smart" {
		if encode {
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
	fmt.Printf(c.Decode())
}

func takeInput() string {
	args := flag.Args()
	stdininfo, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while connecting to stdin: %s\n", err.Error())
	}
	if err == nil && stdininfo.Mode()&os.ModeNamedPipe == 0 {
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
		return os.Args[0]
	}
}
