package ui

import (
	"sync"

	"github.com/empijei/wapty/ui/apis"
)

type UI interface {
	Channel() <-chan apis.Command
}

//TODO Use a map[chan apis.Command]struct{} if nothing else than the channel is used
type UIImpl struct {
	id      int
	channel chan apis.Command
}

func (u *UIImpl) Channel() <-chan apis.Command {
	return u.channel
}

type uis struct {
	sync.RWMutex
	curID int
	list  map[int]UI
}

//This function is used by uis servers to read all the messages from wapty and send them to clients.
func Connect() UI {
	oChans.Lock()
	defer oChans.Unlock()
	toRet := &UIImpl{channel: make(chan apis.Command), id: oChans.curID}
	oChans.list[oChans.curID] = toRet
	oChans.curID++
	return toRet
}

func Disconnect(u UI) {
	oChans.Lock()
	defer oChans.Unlock()
	if u, ok := u.(*UIImpl); ok {
		delete(oChans.list, u.id)
	}
}
