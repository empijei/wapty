package main

import (
	"log"

	"github.com/empijei/Wapty/intercept"
	"github.com/empijei/Wapty/ui"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go ui.MainLoop()
	go ui.ControllerMainLoop()
	intercept.MainLoop()
}
