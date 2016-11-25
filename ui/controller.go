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
var oChans chans

type chans struct {
	sync.RWMutex
	chans []chan Command
}

func init() {
	subScriptions = make(map[string]SubsChannel)
	iChan = make(chan Command)
}

type Subscription struct {
	id          int64
	channel     string
	dataCh      chan Command
	DataChannel <-chan Command
}

func Subscribe(channel string) *Subscription {
	subsMutex.Lock()
	subsCounter++
	//Unless you are sure the out channel will be constantly read, it is strongly
	//suggested to create a buffered channel
	pipe := make(chan Command, 20) //TODO this is arbitrary, give a meaning to this number
	out := Subscription{id: subsCounter, dataCh: pipe, channel: channel}
	if subScriptions[channel] == nil {
		subScriptions[channel] = make(map[int64]Subscription)
	}
	subScriptions[channel][subsCounter] = out
	out.DataChannel = pipe
	subsMutex.Unlock()
	return &out
}

func (s *Subscription) Send(c Command) {
	c.Channel = s.channel
	send(c)
}

//TODO delete this? Dangerous and never used
func UnSubscribe(s *Subscription) {
	subsMutex.RLock()
	defer subsMutex.RUnlock()
	for _, channelSubs := range subScriptions {
		sub, ok := channelSubs[s.id]
		if ok {
			subsMutex.Lock()
			close(sub.dataCh)
			delete(channelSubs, s.id)
			subsMutex.Unlock()
			return
		}
	}
	log.Println("Subscription not found")
}

func send(c Command) {
	oChans.RLock()
	for _, oChan := range oChans.chans {
		oChan <- c
	}
	oChans.RUnlock()
}

//This should be a server's method
func Receive(c Command) {
	iChan <- c
}

//This function is used by uis servers to read all the messages from wapty and send them to clients.
func ConnectUI() <-chan Command {
	toRet := make(chan Command)
	oChans.Lock()
	defer oChans.Unlock()
	oChans.chans = append(oChans.chans, toRet)
	return toRet
}

func MainLoop() {
	for cmd := range iChan {
		subsMutex.RLock()
		for _, out := range subScriptions[cmd.Channel] {
			out.dataCh <- cmd
		}
		subsMutex.RUnlock()
	}
}
