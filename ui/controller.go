//ui is a general high level representation of all the uis connected to the current
//instance of Wapty. Use this from other packages to read user input and write
//output
package ui

import (
	"log"
	"sync"
)

type Command struct {
	Channel string
	Args    Args
	Payload *[]byte
}

type Args map[string]string

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
	oChan <- c
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
