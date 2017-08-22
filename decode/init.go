package decode

import "github.com/empijei/wapty/cli"

var CmdDecode = &cli.Command{
	Name:      "decode",
	Run:       MainStandalone,
	UsageLine: "decode [flags]",
	Short:     "decode something.",
	Long: `decode something in a really clever way:

blah blah blah
	`,
}

var flagEncode bool      // -encode
var flagCodeclist string // -codec

func init() {
	CmdDecode.Flag.BoolVar(&flagEncode, "encode", false, "Sets the decoder to an encoder instead")
	CmdDecode.Flag.StringVar(&flagCodeclist, "codec", "smart",
		`Sets the decoder/encoder codec. Multiple codecs can be specified and comma separated:
	they will be applied one on the output of the previous as in a pipeline.
	`)
}
