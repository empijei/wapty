package intercept

import (
	"reflect"
	"testing"

	"github.com/empijei/wapty/ui/apis"
)

var loopTests = []struct {
	//TODO
}{}

var handleTests = []struct {
	in  apis.Command
	out apis.Command
}{
	{
		apis.Command{
			Action: apis.STN_INTERCEPT,
			Args:   map[apis.ArgName]string{apis.ARG_ON: apis.ARG_FALSE},
		},
		apis.Command{
			Action: apis.STN_INTERCEPT,
			Args:   map[apis.ArgName]string{apis.ARG_ON: apis.ARG_FALSE},
		},
	},
	{
		apis.Command{
			Action: apis.STN_INTERCEPT,
		},
		apis.Command{
			Action: apis.STN_INTERCEPT,
			Args:   map[apis.ArgName]string{apis.ARG_ON: apis.ARG_FALSE},
		},
	},
	{
		apis.Command{
			Action: apis.STN_INTERCEPT,
			Args:   map[apis.ArgName]string{apis.ARG_ON: apis.ARG_TRUE},
		},
		apis.Command{
			Action: apis.STN_INTERCEPT,
			Args:   map[apis.ArgName]string{apis.ARG_ON: apis.ARG_TRUE},
		},
	},
	{
		apis.Command{
			Action: apis.STN_INTERCEPT},
		apis.Command{
			Action: apis.STN_INTERCEPT,
			Args:   map[apis.ArgName]string{apis.ARG_ON: apis.ARG_TRUE},
		},
	},
}

func TestHandleIntercept(t *testing.T) {
	for _, tt := range handleTests {
		out := handleIntercept(tt.in)
		if !reflect.DeepEqual(*out, tt.out) {
			t.Errorf("handleIntercept(%v) => %v, want %v", tt.in, out, tt.out)
		}
	}
}
