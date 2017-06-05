package repeat

import (
	"encoding/json"
	"io"
	"sync"
)

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

func (h *Repeaters) Save(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(h)
}
