package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

//TODO? give a way to run in "forensics" mode and do not touch disk
func init() {
	err := os.MkdirAll(ConfDir, 0700)
	if err != nil {
		panic(err)
	}
}

var (
	// ConfDir holds the path for wapty configuration directory
	ConfDir = path.Join(getUserHomeDir(), ".wapty")
	// ConfName holds the basename of the configuration file
	ConfName = "wapty-conf.json"
)

var conf Configuration

// Configuration represents the struct holding all of wapty settings but not
// the saved status
type Configuration struct {
	RecentProjects []string
	Workspaces     []string
}

func getConfPath() string {
	return ConfDir + string(os.PathSeparator) + ConfName
}

// LoadConf loads the configuration from the default conf file
// WARNING: this function may panic
func LoadConf() {
	buf, err := ioutil.ReadFile(getConfPath())
	switch {
	case os.IsNotExist(err):
	case err != nil:
		panic(err)
	}
	err = json.Unmarshal(buf, &conf)
	if err != nil {
		panic(err)
	}
}

// SaveConf saves the current configuration into the default conf file
// WARNING: this function may panic
func SaveConf() {
	buf, err := json.Marshal(&conf)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(getConfPath(), buf, 0660)
	if err != nil {
		panic(err)
	}
}

func getUserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
