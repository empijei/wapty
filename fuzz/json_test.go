package fuzz

import "testing"

var nestedTests = []struct {
	l   int
	n   string
	exp string
}{
	{
		1,
		"a",
		"{\"a\":{}}",
	},
}

func TestNestedJSON(t *testing.T) {
	for _, tt := range nestedTests {
		out := nestedJSON(tt.l, tt.n)
		if out != tt.exp {
			t.Errorf("nestedJSON(%d,%s) expected %s but got %s", tt.l, tt.n, tt.exp, out)
		}
	}
}
