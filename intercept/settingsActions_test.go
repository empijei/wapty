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
			Action: apis.INTERCEPT,
			Args:   map[apis.Param]string{apis.ON: apis.FALSE},
		},
		apis.Command{
			Action: apis.INTERCEPT,
			Args:   map[apis.Param]string{apis.ON: apis.FALSE},
		},
	},
	{
		apis.Command{
			Action: apis.INTERCEPT,
		},
		apis.Command{
			Action: apis.INTERCEPT,
			Args:   map[apis.Param]string{apis.ON: apis.FALSE},
		},
	},
	{
		apis.Command{
			Action: apis.INTERCEPT,
			Args:   map[apis.Param]string{apis.ON: apis.TRUE},
		},
		apis.Command{
			Action: apis.INTERCEPT,
			Args:   map[apis.Param]string{apis.ON: apis.TRUE},
		},
	},
	{
		apis.Command{
			Action: apis.INTERCEPT},
		apis.Command{
			Action: apis.INTERCEPT,
			Args:   map[apis.Param]string{apis.ON: apis.TRUE},
		},
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
