package intercept

import (
	"reflect"
	"testing"

	"github.com/empijei/Wapty/ui"
)

var loopTests = []struct {
	//TODO
}{}

var handleTests = []struct {
	in  ui.Command
	out ui.Command
}{
	{ui.Command{Action: "intercept", Args: []string{"false"}},
		ui.Command{Action: "intercept", Args: []string{"false"}}},
	{ui.Command{Action: "intercept"},
		ui.Command{Action: "intercept", Args: []string{"false"}}},
	{ui.Command{Action: "intercept", Args: []string{"true"}},
		ui.Command{Action: "intercept", Args: []string{"true"}}},
	{ui.Command{Action: "intercept"},
		ui.Command{Action: "intercept", Args: []string{"true"}},
	},
}

func TestHandleIntercept(t *testing.T) {
	for _, tt := range handleTests {
		out := handleIntercept(tt.in)
		if !reflect.DeepEqual(out, tt.out) {
			t.Errorf("handleIntercept(%v) => %v, want %v", tt.in, out, tt.out)
		}
	}
}
