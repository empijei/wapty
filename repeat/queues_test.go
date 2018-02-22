package repeat

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/empijei/wapty/ui/apis"
)

type MockSubscription struct {
	ID        int64
	Channel   string
	DataCh    chan apis.Command
	SentStuff chan apis.Command
}

func (s *MockSubscription) Receive() apis.Command {
	return <-s.DataCh
}
func (s *MockSubscription) RecChannel() <-chan apis.Command {
	return s.DataCh
}
func (s *MockSubscription) Send(c *apis.Command) {
	s.SentStuff <- *c
}

func TestHandler(t *testing.T) {
	//backupstatus := status
	//backupui := uiRepeater
	//FIXME change this test to directly invoke proper handler and do not spawn
	// the repeater loop
	dataCh := make(chan apis.Command)
	mocksub := &MockSubscription{
		DataCh:    dataCh,
		SentStuff: make(chan apis.Command),
	}
	uiRepeater = mocksub
	status = Repeaters{}
	go RepeaterLoop()
	var req *http.Request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
		_, _ = w.Write([]byte("Test response"))
	}))
	defer ts.Close()
	URL, _ := url.Parse(ts.URL)
	assert := func(condition bool, message string, vars ...interface{}) {
		if !condition {
			t.Errorf(message, vars...)
		}
	}
	dataCh <- apis.Command{
		Channel: apis.CHN_REPEAT,
		Action:  apis.RPT_CREATE,
	}

	tmp := <-mocksub.SentStuff
	id := tmp.Args[apis.ARG_ID]
	assert(id == "0", "Expected repeater id 0 but got "+id)

	tmp = apis.Command{
		Channel: apis.CHN_REPEAT,
		Action:  apis.RPT_GO,
		Payload: []byte(`GET / HTTP/1.1
Host: localhost:` + URL.Port() + `
X-Wapty-Test: TestHeader
Connection: close


`)}
	tmp.PackArgs(
		[]apis.ArgName{apis.ARG_ENDPOINT, apis.ARG_TLS, apis.ARG_ID},
		"localhost:"+URL.Port(), "false", id,
	)
	dataCh <- tmp

	tmp = <-mocksub.SentStuff
	assert(bytes.Contains(tmp.Payload, []byte(`Test response`)), "Unexpected response: "+string(tmp.Payload))
	assert(req.Header.Get("X-Wapty-Test") == "TestHeader", "Test header was not successfully set. Expected <TestHeader> but got <%s>", req.Header.Get("X-Wapty-Test"))
	subid := tmp.Args[apis.ARG_SUBID]
	assert(subid == "0", "Expected repeat payload id 0 but got "+subid)

	tmp = apis.Command{
		Channel: apis.CHN_REPEAT,
		Action:  apis.RPT_GET,
	}
	tmp.PackArgs(
		[]apis.ArgName{apis.ARG_ID, apis.ARG_SUBID},
		id, subid,
	)
	dataCh <- tmp

	tmp = <-mocksub.SentStuff
	var histitem Item
	err := json.Unmarshal(tmp.Payload, &histitem)
	if err != nil {
		t.Error("Unexpected error while fetching repeat entry: " + err.Error())
	}
	assert(bytes.Contains(histitem.Response, []byte(`Test response`)), "Unexpected history response: "+string(tmp.Payload))
	assert(bytes.Contains(histitem.Request, []byte(`TestHeader`)), "Unexpected history request: "+string(tmp.Payload))
	subid = tmp.Args[apis.ARG_SUBID]
	assert(subid == "0", "Expected repeat payload id 0 but got "+subid)
}
