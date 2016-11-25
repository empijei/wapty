package intercept

import (
	"os"
	"testing"

	"github.com/empijei/Wapty/ui"
)

var oChan <-chan ui.Command

func setup() {
	oChan = ui.ConnectUI()
}

func shutdown() {}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
