package intercept

import "github.com/empijei/Wapty/ui"

var uiEditor *ui.Subscription

const EDITORCHANNEL = "proxy/intercept/editor"

//Enum for possible user actions
type Action int

const (
	FORWARDED Action = iota
	EDITED
	DROPPED
	RESPPROVIDED
)

var actions = [...]string{
	"edited",
	"dropped",
	"forwarded",
	"respProvided",
}

func (a Action) String() string {
	return actions[a]
}

var invertActions map[string]Action

func parseAction(s string) Action {
	return invertActions[s]
}

func init() {
	invertActions = make(map[string]Action)
	for i := 0; i <= int(RESPPROVIDED); i++ {
		invertActions[Action(i).String()] = Action(i)
	}
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

func editBuffer(p PayloadType, b *[]byte) (*[]byte, Action) {
	args := ui.Args(map[string]string{"type": p.String()})
	ui.Send(ui.Command{Channel: EDITORCHANNEL, Args: args, Payload: b})
	result := uiEditor.Read()
	action := parseAction(result.Args["action"])
	return result.Payload, action
}
