//ui is a general high level representation of all the uis connected to the current
//instance of Wapty. Use this from other packages to read user input and write
//output
package ui

import (
	"log"
	"sync"
)

//String used to detect the main action within an ui.Args
const ACTION = "action"

type Command struct {
	Channel string
	Action  string
	Args    []string
	Payload []byte
}

//TODO create an helper to build these

type SubsChannel map[int64]Subscription

var subScriptions map[string]SubsChannel
var subsMutex sync.RWMutex
var subsCounter int64
var iChan chan Command
var oChan chan Command

func init() {
	subScriptions = make(map[string]SubsChannel)
	iChan = make(chan Command)
	oChan = make(chan Command)
}

type Subscription struct {
	id      int64
	channel chan Command
}

func (s *Subscription) Read() Command {
	return <-s.channel
}

func Subscribe(channel string) *Subscription {
	subsMutex.Lock()
	subsCounter++
	//Unless you are sure the out channel will be constantly read, it is strongly
	//suggested to create a buffered channel
	pipe := make(chan Command, 20) //TODO this is arbitrary, give a meaning to this number
	out := Subscription{id: subsCounter, channel: pipe}
	if subScriptions[channel] == nil {
		subScriptions[channel] = make(map[int64]Subscription)
	}
	subScriptions[channel][subsCounter] = out
	subsMutex.Unlock()
	return &out
}

//TODO delete this? Dangerous and never used
func UnSubscribe(s *Subscription) {
	subsMutex.RLock()
	defer subsMutex.RUnlock()
	for _, channelSubs := range subScriptions {
		sub, ok := channelSubs[s.id]
		if ok {
			subsMutex.Lock()
			close(sub.channel)
			delete(channelSubs, s.id)
			subsMutex.Unlock()
			return
		}
	}
	log.Println("Subscription not found")
}

func Send(c Command) {
	log.Println("controller sending: ", c)
	oChan <- c
}

func Receive(c Command) {
	iChan <- c
}

//This function is used by uis servers to read all the messages from wapty.
func ConnectUI() <-chan Command {
	//TODO actually duplicate the streams, use the Send function to write to all
	//uis available
	return oChan
}

func MainLoop() {
	for cmd := range iChan {
		subsMutex.RLock()
		for _, out := range subScriptions[cmd.Channel] {
			out.channel <- cmd
		}
		subsMutex.RUnlock()
	}
}
