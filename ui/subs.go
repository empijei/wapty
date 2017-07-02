package ui

import (
	"sync"

	"github.com/empijei/wapty/ui/apis"
)

const SUBBUFFERSIZE = 50

type subsChannel map[int]subscriptionImpl

var subscriptions = make(map[apis.UIChannel]subsChannel)
var subsMutex sync.RWMutex
var subsCounter int

// Subscription is the high-level representation of a connection between a wapty
// component and wapty UI. It will multiplex on apis.UIChannel transparently.
type Subscription interface {
	Receive() apis.Command
	RecChannel() <-chan apis.Command
	Send(*apis.Command)
}

type subscriptionImpl struct {
	id      int
	channel apis.UIChannel
	dataCh  chan apis.Command
}

// Subscribe allows a package to start receiving and sending commands over a apis.UIChannel
func Subscribe(channel apis.UIChannel) Subscription {
	subsMutex.Lock()
	subsCounter++
	//Unless you are sure the out channel will be constantly read, it is strongly
	//suggested to create a buffered channel
	pipe := make(chan apis.Command, SUBBUFFERSIZE)
	out := subscriptionImpl{id: subsCounter, dataCh: pipe, channel: channel}
	if subscriptions[channel] == nil {
		subscriptions[channel] = make(map[int]subscriptionImpl)
	}
	subscriptions[channel][subsCounter] = out
	out.dataCh = pipe
	subsMutex.Unlock()
	return &out
}

// Receive blocks until a command is received
func (s *subscriptionImpl) Receive() apis.Command {
	return <-s.dataCh
}

// RecChannel returns a read-only channel to receive commands. Use this only for
// select statements. If you just need to receiv a command in a blocking way
// please use Receive instead
func (s *subscriptionImpl) RecChannel() <-chan apis.Command {
	return s.dataCh
}

// Send sends the command and sets the channel with the value set in the subscription
func (s *subscriptionImpl) Send(c *apis.Command) {
	c.Channel = s.channel
	send(c)
}
