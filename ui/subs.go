package ui

import "github.com/empijei/wapty/ui/apis"

type Subscription interface {
	Receive() apis.Command
	RecChannel() <-chan apis.Command
	Send(apis.Command)
}

type SubscriptionImpl struct {
	id      int64
	channel string
	dataCh  chan apis.Command
}

func Subscribe(channel string) Subscription {
	subsMutex.Lock()
	subsCounter++
	//Unless you are sure the out channel will be constantly read, it is strongly
	//suggested to create a buffered channel
	pipe := make(chan apis.Command, 50) //TODO this is arbitrary, give a meaning to this number
	out := SubscriptionImpl{id: subsCounter, dataCh: pipe, channel: channel}
	if subScriptions[channel] == nil {
		subScriptions[channel] = make(map[int64]SubscriptionImpl)
	}
	subScriptions[channel][subsCounter] = out
	out.dataCh = pipe
	subsMutex.Unlock()
	return &out
}

func (s *SubscriptionImpl) Receive() apis.Command {
	return <-s.dataCh
}
func (s *SubscriptionImpl) RecChannel() <-chan apis.Command {
	return s.dataCh
}

//Sends the command and sets the channel with the value set in the subscription
func (s *SubscriptionImpl) Send(c apis.Command) {
	c.Channel = s.channel
	send(c)
}
