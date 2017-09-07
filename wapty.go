package main

import (

	//Packages imported for initialization purposes
	_ "github.com/empijei/wapty/decode"
	_ "github.com/empijei/wapty/mocksy"

	"github.com/empijei/wapty/cli"
	"github.com/empijei/wapty/intercept"
	"github.com/empijei/wapty/ui"
)

var cmdProxy = &cli.Cmd{
	Name: "proxy",
	Run: func(...string) {
		go ui.MainLoop()
		intercept.MainLoop()
	},
	UsageLine: "proxy",
	Short:     "work as a proxy",
	Long:      "",
}

func init() {
	cli.AddCommand(cmdProxy)
}

func main() {
	cli.DefaultCommand = cmdProxy
	cli.Init()
}
