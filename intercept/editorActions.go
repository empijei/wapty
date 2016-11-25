package intercept

import (
	"log"

	"github.com/empijei/Wapty/ui"
)

var uiEditor *ui.Subscription

//String used to recognize commands directed to this module
const EDITORCHANNEL = "proxy/intercept/editor"

//Enum for possible user actions
type EditorAction int

const (
	FORWARD EditorAction = iota
	EDIT
	DROP
	PROVIDERESP
)

var editorActions = [...]string{
	"forward",
	"edit",
	"drop",
	"provideResp",
}

func (a EditorAction) String() string {
	return editorActions[a]
}

var invertEditorActions map[string]EditorAction

func parseEditorAction(s string) EditorAction {
	return invertEditorActions[s]
}

//Enum for possible payloads types
type PayloadType int

const (
	REQUEST PayloadType = iota
	RESPONSE
)

var payloads = [...]string{
	"request",
	"response",
}

func (p PayloadType) String() string {
	return payloads[p]
}

func init() {
	invertEditorActions = make(map[string]EditorAction)
	for i := 0; i <= int(PROVIDERESP); i++ {
		invertEditorActions[EditorAction(i).String()] = EditorAction(i)
	}
	uiEditor = ui.Subscribe(EDITORCHANNEL)
	uiHistory = ui.Subscribe(HISTORYCHANNEL)
}

func editBuffer(p PayloadType, b []byte) ([]byte, EditorAction) {
	log.Println("Editing " + p.String())
	args := []string{p.String()}
	ui.Send(ui.Command{Channel: EDITORCHANNEL, Action: "edit", Args: args, Payload: b}) //TODO add Action?
	log.Println("Waiting for user interaction")
	result := <-uiEditor.Channel
	log.Println("User interacted")
	//FIXME do something if action not recognized!
	action := parseEditorAction(result.Action)
	return result.Payload, action
}
