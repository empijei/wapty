package intercept

import (
	"bytes"
	"testing"

	"github.com/empijei/wapty/ui/apis"
)

type paramsType struct {
	p   apis.Action
	cmd apis.Command
	b   []byte
}

type outputType struct {
	b []byte
	e apis.Action
}

var editBufferTests = []struct {
	in  paramsType
	out outputType
}{
//TODO
}

func TestEditBuffer(t *testing.T) {
	mockChan := make(chan apis.Command)
	uiEditor = &MockSubscription{DataCh: mockChan}
	//uiEditor = &ui.SubscriptionImpl{
	//dataCh: mockChan,
	//}
	defer func() {
		uiEditor = nil
		close(mockChan)
	}()
	for i, tt := range editBufferTests {
		go func() {
			mockChan <- tt.in.cmd
		}()
		b, e := editBuffer(tt.in.p, tt.in.b, "https://thisisatest.com:443")
		if bytes.Compare(b, tt.out.b) != 0 || e != tt.out.e {
			t.Errorf("editBufferTests[%d]", i)
		}
	}
}
