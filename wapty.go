package main

import "github.com/empijei/Wapty/intercept"
import "github.com/empijei/Wapty/ui"

func main() {
	go ui.MainLoop()
	intercept.MainLoop()
}
