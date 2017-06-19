package intercept

import "github.com/empijei/wapty/ui/apis"

type MockSubscription struct {
	ID        int64
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

//Sends the command and sets the channel with the value set in the subscription
func (s *MockSubscription) Send(c apis.Command) {
	s.SentStuff = append(s.SentStuff, c)
}
