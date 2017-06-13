package repeat

import (
	"testing"

	"github.com/empijei/wapty/ui"
	"github.com/empijei/wapty/ui/apis"
)

type MockSubscription struct {
	Id        int64
	Channel   string
	DataCh    chan apis.Command
	SentStuff []apis.Command
}

func (s *MockSubscription) Receive() apis.Command {
	return <-s.DataCh
}
func (s *MockSubscription) RecChannel() <-chan apis.Command {
	return s.DataCh
}

func (s *MockSubscription) Send(c apis.Command) {
	s.SentStuff = append(s.SentStuff, c)
}

func TestHandleGo(t *testing.T) {
	uirbak := uiRepeater
	defer func() {
		uiRepeater = uirbak
	}()
	testSubc := make(chan apis.Command)
	testSub := &MockSubscription{
		DataChannel: testchan,
	}
	testUI := ui.UI{}
}
