package main

import (
	"github.com/empijei/Wapty/intercept"
	"github.com/empijei/Wapty/ui"
	"github.com/empijei/Wapty/ui/websock"
)

func main() {
	go ui.MainLoop()
	go websock.MainLoop()
	intercept.MainLoop()
}
