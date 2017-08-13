package help

import "github.com/empijei/wapty/common"

var CmdHelp = &common.Command{
	Name:      "help",
	Run:       Main,
	UsageLine: "help",
	Short:     "display help information for wapty commands",
	Long:      "",
}
