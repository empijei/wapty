package repeat

import (
	"encoding/json"
	"io"
	"sync"
)

var status Repeaters

// Repeaters is a list of repeater with its embedded RWMutex
type Repeaters struct {
	sync.RWMutex
	Repeats []*Repeater
}

// Add appends in a thread-safe way a repeater to the current status
func (h *Repeaters) Add(r *Repeater) {
	h.Lock()
	defer h.Unlock()
	h.Repeats = append(h.Repeats, r)
}

// Save writes the current status to the given writer in a thread-safe way
func (h *Repeaters) Save(w io.Writer) error {
	h.RLock()
	defer h.RUnlock()
	e := json.NewEncoder(w)
	return e.Encode(h)
}
