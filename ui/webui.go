package ui

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	rice "github.com/GeertJohan/go.rice"
)

var tabs = []struct {
	Name  string
	Title string
}{
	{
		"proxy",
		"Proxy",
	},
	{
		"history",
		"HTTP History",
	},
}

type Tab struct {
	Title   string
	Active  bool
	ID      string
	Content string
}

type Index struct {
	Tabs []*Tab
}

var names = make(map[string]struct{})

var appPage []byte

var templatesFolder *rice.Box

func loadTemplates() {
	templatesFolder = rice.MustFindBox("templates")
	indexraw := mustreadall("index.tmpl")
	indextmpl := template.Must(template.New("index").Parse(indexraw))
	home := new(Index)
	for i, tab := range tabs {
		home.Tabs = append(home.Tabs, loadTab(tab.Name, tab.Title, i == 0))
	}
	buf := bytes.NewBuffer(nil)
	err := indextmpl.Execute(buf, home)
	if err != nil {
		panic(err)
	}
	appPage = buf.Bytes()
}

func loadTab(name string, title string, active bool) *Tab {
	//names must be unique
	if _, ok := names[name]; ok {
		panic(fmt.Errorf("name %s already in use!", name))
	}
	names[name] = struct{}{}

	return &Tab{
		Title:   title,
		Active:  active,
		ID:      name,
		Content: mustreadall(name + ".html"),
	}
}

func mustreadall(path string) string {
	file, err := templatesFolder.Open(path)
	if err != nil {
		panic(err)
	}
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(buf)
}
