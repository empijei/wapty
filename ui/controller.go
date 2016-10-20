package ui

import (
	"log"
	"sync"
)

//only listen for localhost
type Command struct {
	Channel string
	Args    map[string]string
}

type Args map[string]string

var subScriptions map[string]map[int64]Subscription
var subsMutex sync.RWMutex
var subsCounter int64
var ioChan chan Command

func init() {
	subScriptions = make(map[string]map[int64]Subscription)
	ioChan = make(chan Command)
}

type Subscription struct {
	id      int64
	channel chan Command
}

func (s *Subscription) Read() Command {
	return <-s.channel
}

//Unless you are sure the out channel will be constantly read, it is strongly
//suggested to create a buffered channel
func SubScribe(channel string) Subscription {
	subsMutex.Lock()
	subsCounter++
	pipe := make(chan Command, 20) //TODO this is arbitrary, give a meaning to this number
	out := Subscription{id: subsCounter, channel: pipe}
	if subScriptions[channel] == nil {
		subScriptions[channel] = make(map[int64]Subscription)
	}
	subScriptions[channel][subsCounter] = out
	subsMutex.Unlock()
	return out
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
	ioChan <- c
}

func MainLoop() {
	for cmd := range ioChan {
		subsMutex.RLock()
		for _, out := range subScriptions[cmd.Channel] {
			out.channel <- cmd
		}
		subsMutex.RUnlock()
	}
}
