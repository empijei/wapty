package repeat

import (
	"github.com/empijei/Wapty/ui"
)

var uiRepeater *ui.Subscription

//String used to recognize commands directed to this module
const REPEATCHANNEL = "repeat"

//Enum for possible user actions
type RepeaterAction int

const (
	//Creates a new repeater entry
	CREATE RepeaterAction = iota
	//Performs the request
	GO
	//Retrieves an history item
	GET
)

var repeaterActions = [...]string{
	"create",
	"go",
	"get",
}

func (a RepeaterAction) String() string {
	return repeaterActions[a]
}

var invertRepeaterActions map[string]RepeaterAction

func parseRepeaterAction(s string) RepeaterAction {
	return invertRepeaterActions[s]
}

func init() {
	invertRepeaterActions = make(map[string]RepeaterAction)
	for i := 0; i <= 0; i++ {
		invertRepeaterActions[RepeaterAction(i).String()] = RepeaterAction(i)
	}
	uiRepeater = ui.Subscribe(REPEATCHANNEL)
}

/*
func editBuffer(p PayloadType, b []byte, endpoint string) ([]byte, EditorAction) {
	log.Println("Editing " + p.String())
	//result := ui.Command{Action: "edit", Payload: b, Channel: EDITORCHANNEL}
	args := []string{p.String(), endpoint}
	uiEditor.Send(ui.Command{Action: "edit", Args: args, Payload: b})
	log.Println("Waiting for user interaction")
	result := <-uiEditor.DataChannel
	log.Println("User interacted")
	//FIXME do something if action not recognized!
	action := parseEditorAction(result.Action)
	return result.Payload, action
}
*/
