//Package ui is a general high level representation of all the uis connected to the current
//instance of Wapty. Use this from other packages to read user input and write
//output
package ui

import (
	"sync"

	"github.com/empijei/wapty/ui/apis"
)

//String used to detect the main action within an ui.Args
const ACTION = "action"

type SubsChannel map[int64]SubscriptionImpl

var subScriptions map[apis.UIChannel]SubsChannel
var subsMutex sync.RWMutex
var subsCounter int64
var iChan chan apis.Command
var oChans uis

func init() {
	subScriptions = make(map[apis.UIChannel]SubsChannel)
	iChan = make(chan apis.Command)
	oChans.list = make(map[int]UI)
}

func send(c apis.Command) {
	oChans.RLock()
	for _, oChan := range oChans.list {
		if oChan, ok := oChan.(*UIImpl); ok {
			oChan.channel <- c
		}
	}
	oChans.RUnlock()
}

//This should be a server's method
func Receive(c apis.Command) {
	iChan <- c
}

func ControllerMainLoop() {
	for cmd := range iChan {
		subsMutex.RLock()
		for _, out := range subScriptions[cmd.Channel] {
			out.dataCh <- cmd
		}
		subsMutex.RUnlock()
	}
}
