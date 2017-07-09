package decode

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func init() {
	RegisterFlagStandalone()
}

func TestMainStandalone(t *testing.T) {
	if os.Getenv("BE_MAIN") == "1" {
		MainStandalone()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMainStandalone", "-codec", "b16,b64", "5a6d3976")
	cmd.Env = append(os.Environ(), "BE_MAIN=1")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stdout
	err := cmd.Run()

	if err != nil {
		panic(err)
	}

	res := string(out.Bytes())
	if strings.HasPrefix(res, "fooPASS") == false {
		t.Errorf("Expected decoded value: %s, but got '%s'", "foo", res)
	}
}
