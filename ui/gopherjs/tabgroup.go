package main

import js "github.com/gopherjs/gopherjs/js"

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

const tabHeaderstmpl = `<ul id="tab-list" class="nav nav-tabs" role="tablist">
</ul>
`

const tabBodiestmpl = `<div id="tab-content" class="tab-content">
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

	js.Global.Get("document").Call("createElement")
	//tg.tabHeaders.Call("appendChild", th)
	//tg.tabBodies.Call("appendChild", tb)
	return 0
}
