package ui

import "github.com/empijei/wapty/ui/apis"

// Subscription is the high-level representation of a connection between a wapty
// component and wapty UI. It will multiplex on apis.UIChannel transparently.
type Subscription interface {
	Receive() apis.Command
	RecChannel() <-chan apis.Command
	Send(apis.Command)
}

type subscriptionImpl struct {
	id      int64
	channel apis.UIChannel
	dataCh  chan apis.Command
}

// Subscribe allows a package to start receiving and sending commands over a apis.UIChannel
func Subscribe(channel apis.UIChannel) Subscription {
	subsMutex.Lock()
	subsCounter++
	//Unless you are sure the out channel will be constantly read, it is strongly
	//suggested to create a buffered channel
	pipe := make(chan apis.Command, 50) //TODO this is arbitrary, give a meaning to this number
	out := subscriptionImpl{id: subsCounter, dataCh: pipe, channel: channel}
	if subScriptions[channel] == nil {
		subScriptions[channel] = make(map[int64]subscriptionImpl)
	}
	subScriptions[channel][subsCounter] = out
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
func (s *subscriptionImpl) Send(c apis.Command) {
	c.Channel = s.channel
	send(c)
}
