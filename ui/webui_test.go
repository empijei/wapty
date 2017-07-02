package ui

import "testing"

func TestLoad(t *testing.T) {
	loadTemplates()
	t.Log(string(appPage))
}
