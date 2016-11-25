package intercept

import (
	"encoding/json"
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
		ui.Command{Channel: SETTINGSCHANNEL, Action: "intercept", Args: []string{"false"}}},
	{ui.Command{Action: "intercept"},
		ui.Command{Action: "intercept", Channel: SETTINGSCHANNEL, Args: []string{"false"}}},
	{ui.Command{Action: "intercept", Args: []string{"true"}},
		ui.Command{Channel: SETTINGSCHANNEL, Action: "intercept", Args: []string{"true"}}},
	{ui.Command{Action: "intercept"},
		ui.Command{Action: "intercept", Channel: SETTINGSCHANNEL, Args: []string{"true"}},
	},
}

func TestHandleIntercept(t *testing.T) {
	for _, tt := range handleTests {
		out := handleIntercept(tt.in)
		actualout, _ := json.MarshalIndent(out, " ", " ")
		expectedout, _ := json.MarshalIndent(tt.out, " ", " ")
		if string(actualout) != string(expectedout) {
			t.Errorf("handleIntercept(%v) => %v, want %v", tt.in, out, tt.out)
		}
	}
}
