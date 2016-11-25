package intercept

import (
	"bytes"
	"testing"

	"github.com/empijei/Wapty/ui"
)

type paramsType struct {
	p   PayloadType
	cmd ui.Command
	b   []byte
}

type outputType struct {
	b []byte
	e EditorAction
}

var editBufferTests = []struct {
	in  paramsType
	out outputType
}{
//TODO
}

func TestEditBuffer(t *testing.T) {
	mockChan := make(chan ui.Command)
	uiEditor = &ui.Subscription{Channel: mockChan}
	defer func() {
		uiEditor = nil
		close(mockChan)
	}()
	for i, tt := range editBufferTests {
		go func() {
			mockChan <- tt.in.cmd
		}()
		b, e := editBuffer(tt.in.p, tt.in.b)
		if bytes.Compare(b, tt.out.b) != 0 || e != tt.out.e {
			t.Errorf("editBufferTests[%d]", i)
		}
	}
}
