package common

import (
	"flag"
	"fmt"
	"os"
)

// Command is used by any package exposing a runnable command to gather information
// about command name, usage and flagset.
type Command struct {
	Name      string
	Run       func()
	UsageLine string
	Short     string
	Long      string
	Flag      flag.FlagSet
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	c.Flag.PrintDefaults()
	os.Exit(2)
}
