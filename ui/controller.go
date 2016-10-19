package ui

import "sync"

//only listen for localhost
type Command struct {
	Channel string
	Args    map[string]string
}

var subScriptions map[string][]chan<- Command
var subsMutex sync.Mutex
var ioChan chan Command

func init() {
	subScriptions = make(map[string][]chan<- Command)
	ioChan = make(chan Command)
}

//Unless you are sure the out channel will be constantly read, it is strongly
//suggested to create a buffered channel
func SubScribe(channel string, out chan<- Command) {
	subsMutex.Lock()
	defer subsMutex.Unlock()
	subScriptions[channel] = append(subScriptions[channel], out)
}

func Send(c Command) {
	ioChan <- c
}

func MainLoop() {
	for cmd := range ioChan {
		subsMutex.Lock()
		for _, out := range subScriptions[cmd.Channel] {
			out <- cmd
		}
		subsMutex.Unlock()
	}
}
