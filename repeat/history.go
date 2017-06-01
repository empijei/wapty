package repeat

import (
	"bytes"
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

//FIXME
func (h *Repeaters) Save() io.Reader {
	b := bytes.NewBuffer(nil)
	e := json.NewEncoder(b)
	go e.Encode(h)
	return b
}
