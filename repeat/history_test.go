package repeat

import (
	"io/ioutil"
	"testing"
)

func TestSave(t *testing.T) {
	rr := NewRepeater()
	ri := RepeatItem{
		Host:     "host:port",
		Request:  []byte("Request"),
		Response: []byte("Response"),
	}
	rr.history = append(rr.history, ri)
	status.Add(rr)
	r := status.Save()
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(buf))
}
