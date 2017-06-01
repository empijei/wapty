package intercept

import (
	"log"

	"github.com/empijei/wapty/ui"
	"github.com/empijei/wapty/ui/apis"
)

var uiEditor *ui.Subscription

func init() {
	uiEditor = ui.Subscribe(apis.EDITORCHANNEL)
}

//Invokes the edit action on the proxy ui. When a response is received returns
//the payload and the action in its string form. It does not attempt to validate
//the action, the caller must take care of it.
func editBuffer(p apis.PayloadType, b []byte, endpoint string) ([]byte, string) {
	log.Println("Editing " + p.String())
	//result := ui.Command{Action: "edit", Payload: b, Channel: EDITORCHANNEL}
	args := []string{p.String(), endpoint}
	uiEditor.Send(ui.Command{Action: "edit", Args: args, Payload: b})
	log.Println("Waiting for user interaction")
	result := <-uiEditor.DataChannel
	log.Println("User interacted")
	//FIXME do something if action not recognized!
	return result.Payload, result.Action
}
