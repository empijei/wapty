package main

import (
	"log"

	"github.com/empijei/Wapty/intercept"
	"github.com/empijei/Wapty/ui"
	"github.com/empijei/Wapty/ui/websock"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go ui.MainLoop()
	go websock.MainLoop()
	intercept.MainLoop()
}
