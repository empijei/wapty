package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/empijei/Wapty/decode"
	"github.com/empijei/Wapty/intercept"
	"github.com/empijei/Wapty/ui"
)

var mode = flag.String("mode", "proxy", "Selects the mode Wapty should be started on, available values are: proxy, decode")
var encode = flag.Bool("encode", false, "In decode mode sets the decoder to an encoder instead")
var codec = flag.String("codec", "smart", "In decode mode sets the decoder/encoder codec. \n"+
	"Multiple codecs can be specified and comma separated, they will be applied one on the output of the previous.")

func main() {
	flag.Parse()
	switch *mode {
	case "decode":
		fmt.Fprintln(os.Stderr, "Running in decode/encode mode")
		decode.MainStandalone(*codec, *encode)
		return
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go ui.MainLoop()
	go ui.ControllerMainLoop()
	intercept.MainLoop()
}
