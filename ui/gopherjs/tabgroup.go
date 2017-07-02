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

	//Creating the tab header
	li := createElement("li")
	//Creating the link to the tab content
	a := createElement("a")
	a.SetAttributes(map[string]string{
		"href":        "#tab" + tg.Name + strconv.Itoa(tg.currentID),
		"role":        "tab",
		"data-toggle": "tab"})
	a.SetTextContentf("Tab header for tab %d of tabgroup %s", tg.currentID, tg.Name)
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

	//tg.tabHeaders.Call("appendChild", th)
	//tg.tabBodies.Call("appendChild", tb)
	//BOOKMARK
	return tg.currentID
}
