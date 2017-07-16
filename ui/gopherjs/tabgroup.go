// +build js

package main

import "strconv"

const tabtitletmpl = `
<li>
<a href="#{{.GroupName}}tab{{.ID}}" role="tab" data-toggle="tab">
Tab {{.Title}}
<button class="close" type="button" title="Remove this entry">
✗
</button>
</a>
</li>
`

const tabcontenttmpl = `
<div class="tab-pane fade" id="{{.Group}}tab{{.ID}}">
{{.Content}}
</div>
`

const tabHeadersContainertmpl = `<ul id="tab-list" class="nav nav-tabs" role="tablist">
</ul>
`

const tabBodiesContainertmpl = `<div id="tab-content" class="tab-content">
   <div class="tab-pane fade in active" id="tab1">Tab 1 content</div>
`

type TabGroup struct {
	Name       string
	parentNode DomElement
	tabHeaders DomElement
	tabBodies  DomElement
	currentID  int
	Tabs       []Tab
}

func newTabGroup(name string, parentNode DomElement) *TabGroup {
	//BOOKMARK
}

type Tab struct {
	GroupName string
	ID        int
	Title     string
	Content   string
}

func (tg *TabGroup) addTab(title string, content string) (tabID int) {
	tg.currentID++
	t := Tab{
		GroupName: tg.Name,
		ID:        tg.currentID,
		Title:     title,
		Content:   content,
	}
	tg.Tabs = append(tg.Tabs, t)

	refid := "tab" + tg.Name + strconv.Itoa(tg.currentID)
	//Creating the tab header
	li := createElement("li")
	//Creating the link to the tab content
	a := createElement("a")
	a.SetAttributes(map[string]string{
		"href":        "#" + refid,
		"role":        "tab",
		"data-toggle": "tab"})
	a.SetTextContentf("Tab header for tab %d of tabgroup %s, %s", tg.currentID, tg.Name, title)
	//Creating the ✗ button for the tab
	button := createElement("button")
	button.SetAttributes(map[string]string{
		"class": "close",
		"type":  "button",
		"title": "Remove this tab"})
	button.SetTextContent("✗")
	//Putting all together
	a.AppendChild(button)
	li.AppendChild(a)
	tg.tabHeaders.AppendChild(li)

	//Creating tab body
	div := createElement("div")
	div.SetAttributes(map[string]string{
		"class": "tab-pane fade",
		id:      refid,
	})
	div.Set("innerHTML", content)
	//Adding tab body to DOM
	tg.tabBodies.AppendChild(div)

	return tg.currentID
}
