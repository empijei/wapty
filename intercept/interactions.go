package intercept

import (
	"log"

	"github.com/empijei/Wapty/ui"
)

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
	"forwarded",
	"edited",
	"dropped",
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
	log.Println("Editing " + p.String())
	args := ui.Args(map[string]string{"type": p.String()})
	ui.Send(ui.Command{Channel: EDITORCHANNEL, Args: args, Payload: b})
	log.Println("Waiting for user interaction")
	result := uiEditor.Read()
	log.Println("User interacted")
	action := parseAction(result.Args["action"]) //TODO make this a const
	return result.Payload, action
}
