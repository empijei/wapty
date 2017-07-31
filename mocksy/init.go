package mocksy

import (
	"io"
	"os"
)

var outw io.Writer

func init() {
	responseHistory = make([]Item, 0)
	outw = os.Stderr
}
