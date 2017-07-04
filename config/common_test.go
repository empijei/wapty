package config

import (
	"os"
	"reflect"
	"testing"
)

var (
	tmpConfDir string
	savedConf  Configuration
)

func setup() {
	tmpConfDir = ConfDir
	ConfDir = os.TempDir()
	savedConf = conf
}

func shutdown() {
	ConfDir = tmpConfDir
	conf = savedConf
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestSaveLoadConf(t *testing.T) {
	conf = Configuration{
		RecentProjects: []string{"foo", "bar"},
		Workspaces:     []string{"lol", "lal"},
	}
	compare := conf
	SaveConf()
	conf = *new(Configuration)
	LoadConf()
	if !reflect.DeepEqual(compare, conf) {
		t.Errorf("Failed Save or Load: expected %#v but got %#v", compare, conf)
	}
}
