package repeat

import (
	"testing"

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

type MockUI struct {
	channel chan apis.Command
}

func (u *MockUI) Channel() <-chan apis.Command {
	return u.channel
}

func TestHandleGo(t *testing.T) {
	/*
		uirbak := uiRepeater
		defer func() {
			uiRepeater = uirbak
		}()
		testSubCh := make(chan apis.Command)
		testSub := &MockSubscription{
			DataCh: testSubCh,
		}
		uiRepeater = testSub

		testChan := make(chan RepTest, 2)
		input := make(chan []byte, 2)
		var l net.Listener
		go listener(t, testChan, input, &l)
		defer func() { _ = l.Close() }()
		//BOOKMARK
		for _, tt := range RepeatTests {
			testChan <- tt
			defaultTimeout = 1 * time.Second
			cmd := &apis.Command{
				Action:apis.GET,
				Channel:apis.REPEATCHANNEL,
				Payload:tt.in,
				Args:[]string{"","",""}
			}
			r := NewRepeater()
			buf := bytes.NewBuffer(tt.in)
			res, err := r.Repeat(buf, "localhost:12321", false)
			if err != nil {
				t.Error(err)
				return
			}
			resBuf, err := ioutil.ReadAll(res)
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(resBuf, tt.out) != 0 {
				t.Errorf("Expected <%s> but got <%s>", string(tt.out), string(resBuf))
			}
			in := <-input
			if bytes.Compare(in, tt.in) != 0 {
				t.Errorf("Expected <%s> but got <%s>", string(tt.in), string(in))
			}
		}
	*/
}
