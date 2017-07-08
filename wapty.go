package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/empijei/wapty/decode"
	"github.com/empijei/wapty/intercept"
	"github.com/empijei/wapty/mocksy"
	"github.com/empijei/wapty/ui"
)

var (
	//Version is taken by the build flags, represent current version as
	//<major>.<minor>.<patch>
	Version string

	//Commit is the output of `git rev-parse HEAD` at the moment of the build
	Commit string
)

var commands = []struct {
	name string
	main func()
}{
	{"decode", decode.MainStandalone},
	{"proxy", proxyMain},
	{"mocksy", mocksy.Main},
	{"version", func() {
		// Setup fallback version and commit in case wapty wasn't "properly" compiled
		if len(Version) == 0 {
			Version = "Unknown"
		}
		if len(Commit) == 0 {
			Commit = "Unknown"
		}
		fmt.Printf("Version: %s\nCommit: %s\n", Version, Commit)
	}},
}

func init() {
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	if len(os.Args) > 1 {
		//read the first argument
		directive := os.Args[1]
		if len(os.Args) > 2 {
			//shift parameters left, but keep argv[0]
			os.Args = append(os.Args[:1], os.Args[2:]...)
		} else {
			os.Args = os.Args[:1]
		}
		invokeMain(directive)
	} else {
		proxyMain()
	}
}

func proxyMain() {
	go ui.MainLoop()
	intercept.MainLoop()
}

func invokeMain(s string) {
	var toinvoke func()
	var success bool
	for _, cmd := range commands {
		if cmd.name == s {
			toinvoke = cmd.main
			success = true
			break
		}
		if strings.HasPrefix(cmd.name, s) {
			if toinvoke != nil {
				fmt.Fprintf(os.Stderr, "Ambiguous command: '%s'. ", s)
				success = false
			} else {
				toinvoke = cmd.main
				success = true
			}
		}
	}
	if success {
		toinvoke()
		return
	}
	if toinvoke == nil {
		fmt.Fprintf(os.Stderr, "Command not found: '%s'. ", s)
	}

	fmt.Fprintln(os.Stderr, "Available commands are: ")
	for _, cmd := range commands {
		fmt.Fprintln(os.Stderr, "\t"+cmd.name)
	}
}
