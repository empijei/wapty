package cli

var CmdHelp = &Command{
	Name:      "help",
	Run:       Main,
	UsageLine: "help",
	Short:     "display help information for wapty commands",
	Long:      "",
}
