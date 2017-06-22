package intercept

import (
	"log"

	"github.com/empijei/wapty/ui"
	"github.com/empijei/wapty/ui/apis"
)

var uiEditor ui.Subscription

func init() {
	uiEditor = ui.Subscribe(apis.EDITORCHANNEL)
}

//Invokes the edit action on the proxy ui. When a response is received returns
//the payload and the action in its string form. It does not attempt to validate
//the action, the caller must take care of it.
func editBuffer(p apis.Action, b []byte, endpoint string) ([]byte, apis.Action) {
	log.Println("Editing " + p)
	args := map[apis.Param]string{
		apis.PAYLOADTYPE: string(p),
		apis.ENDPOINT:    endpoint}
	uiEditor.Send(apis.Command{Action: apis.EDIT, Args: args, Payload: b})
	log.Println("Waiting for user interaction")
	result := uiEditor.Receive()
	log.Println("User interacted")
	//FIXME do something if action not recognized!
	return result.Payload, result.Action
}
