package main

import "github.com/empijei/wapty/ui/apis"

func handleRepeat(mdg apis.Command) {

}

//DOM Stuff
const tabtitletmpl = `
<li>
<a href="#reptab{{.ID}}" role="tab" data-toggle="tab">
Tab {{.Title}}
<button class="close" type="button" title="Remove this entry">
✗
</button>
</a>
</li>
`

const tabcontenttmpl = `
<div class="tab-pane fade" id="reptab{{.ID}}">
Tab {{.tabID}} content
</div>
`

func addTab(title, request, response string) {}
