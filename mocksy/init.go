package mocksy

import (
	"io"
	"os"

	"github.com/empijei/wapty/cli"
)

var outw io.Writer

var CmdMocksy = &cli.Command{
	Name:      "mocksy",
	Run:       Main,
	UsageLine: "mocksy",
	Short:     "mock responses from a server",
	Long:      "",
}

func init() {
	responseHistory = make([]Item, 0)
	outw = os.Stderr
}
