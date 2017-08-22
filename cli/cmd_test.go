package cli

import "testing"

var findCommandTests = []struct {
	in  string
	out string
}{
	{"pref", "pref"},
	{"prefixs", "prefixsuffix"},
	{"prefix", ""},
	{"suffixes", ""},
	{"prefixd", "prefixdiffsuffix"},
}

func TestFindCommand(t *testing.T) {
	bk := WaptyCommands
	defer func() { WaptyCommands = bk }()

	var out string
	WaptyCommands = []*Cmd{
		{
			Name: "prefixsuffix",
			Run: func(_ ...string) {
				out = "prefixsuffix"
			},
		},
		{
			Name: "prefixdiffsuffix",
			Run: func(_ ...string) {
				out = "prefixdiffsuffix"
			},
		},
		{
			Name: "pref",
			Run: func(_ ...string) {
				out = "pref"
			},
		},
		{
			Name: "suffix",
			Run: func(_ ...string) {
				out = "suffix"
			},
		},
	}
	for _, tt := range findCommandTests {
		out = ""
		c, e := FindCommand(tt.in)
		if e != nil {
			if tt.out != "" {
				t.Errorf("Expected command <%s> but got error: <%s> instead", tt.out, e.Error())
				continue
			}
			continue
		}
		c.Run()
		if tt.out != out {
			t.Errorf("Expected <%s> but got <%s>", tt.out, out)
		}
	}
}
