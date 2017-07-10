package decode

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var CodecPayloadTest = []struct {
	cFlag   string
	in      string
	eOut    string
	success bool
}{
	{
		"b16",
		"666F6F626172",
		"foobar",
		true,
	},
	{
		"b16,b64",
		"5a6d3976",
		"foo",
		true,
	},
	{
		",,",
		"5a6d3976",
		"",
		false,
	},
	{
		"",
		"5a6d3976",
		"",
		false,
	},
	{
		"smart,smart",
		"5a6d3976",
		"foo",
		true,
	},
}

func init() {
	RegisterFlagStandalone()
}

func TestMainStandalone(t *testing.T) {
	if os.Getenv("BE_MAIN") == "1" {
		MainStandalone()
		return
	}

	for _, tt := range CodecPayloadTest {
		cmd := exec.Command(os.Args[0], "-test.run=TestMainStandalone", "-codec", tt.cFlag, tt.in)
		cmd.Env = append(os.Environ(), "BE_MAIN=1")

		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()

		res := string(out.Bytes())

		if e, ok := err.(*exec.ExitError); ok && e.Success() != tt.success {
			t.Errorf("Expected exit code: %v, but got %v", tt.success, e.Success())
		}

		if !strings.HasPrefix(res, tt.eOut) {
			t.Errorf("Expected decoded value: %s, but got '%s'", tt.eOut, res)
		}
	}
}
