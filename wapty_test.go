package main

import "testing"

var InvokeMainTests = []struct {
	in  string
	out string
}{
	{"pref", "pref"},
	{"prefixs", "prefixsuffix"},
	{"prefix", ""},
	{"suffixes", ""},
	{"prefixd", "prefixdiffsuffix"},
}

func TestInvokeMain(t *testing.T) {
	var out string
	commands = []struct {
		name string
		main func()
	}{
		{
			"prefixsuffix",
			func() {
				out = "prefixsuffix"
			},
		},
		{
			"prefixdiffsuffix",
			func() {
				out = "prefixdiffsuffix"
			},
		},
		{
			"pref",
			func() {
				out = "pref"
			},
		},
		{
			"suffix",
			func() {
				out = "suffix"
			},
		},
	}
	for _, tt := range InvokeMainTests {
		out = ""
		invokeMain(tt.in)
		if tt.out != out {
			t.Errorf("Expected %s but got %s", tt.out, out)
		}
	}
}
