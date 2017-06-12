package repeat

import (
	"bytes"
	"testing"
)

func TestSave(t *testing.T) {
	rr := NewRepeater()
	ri := RepeatItem{
		Host:     "host:port",
		Request:  []byte("Request"),
		Response: []byte("Response"),
	}
	rr.History = append(rr.History, ri)
	status.Add(rr)
	b := bytes.NewBuffer(nil)
	err := status.Save(b)
	if err != nil {
		t.Log(err)
	}
	//FIXME
	t.Log(string(b.Bytes()))
}
