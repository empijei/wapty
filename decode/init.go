package decode

import "github.com/empijei/wapty/cli"

var cmdDecode = &cli.Cmd{
	Name:      "decode",
	Run:       MainStandalone,
	UsageLine: "decode [flags]",
	Short:     "decode something.",
	//FIXME write this, and add that this can read from pipe
	Long: `decode something in a really clever way:

blah blah blah
	`,
}

var flagEncode bool      // -encode
var flagCodeclist string // -codec

func init() {
	cmdDecode.Flag.BoolVar(&flagEncode, "encode", false, "Sets the decoder to an encoder instead")
	cmdDecode.Flag.StringVar(&flagCodeclist, "codec", "smart",
		`Sets the decoder/encoder codec. Multiple codecs can be specified and comma separated:
	they will be applied one on the output of the previous as in a pipeline.
	`)
	cli.AddCommand(cmdDecode)
}
