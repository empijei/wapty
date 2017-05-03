package repeat

import "sync"

var status Repeaters

type Repeaters struct {
	sync.RWMutex
	Repeats []*Repeater
}

func (h *Repeaters) Add(r *Repeater) {
	h.Lock()
	defer h.Unlock()
	h.Repeats = append(h.Repeats, r)
}
