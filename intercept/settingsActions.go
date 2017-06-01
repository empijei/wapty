package intercept

import (
	"log"

	"github.com/empijei/wapty/ui"
	"github.com/empijei/wapty/ui/apis"
)

func settingsLoop() {
	for {
		select {
		case cmd := <-uiSettings.DataChannel:
			log.Println("Settings accessed")
			switch cmd.Action {
			case apis.INTERCEPT.String():
				uiSettings.Send(handleIntercept(cmd))
			default:
				//TODO send error?
				log.Printf("Unknown action: %v\n", cmd.Action)
			}
		case <-done:
			return
		}
	}
}

func handleIntercept(cmd ui.Command) ui.Command {
	if len(cmd.Args) >= 1 {
		log.Println("Requested change intercept status")
		intercept.SetValue(cmd.Args[0] == "true" || cmd.Args[0] == "on")
		value := "false"
		if intercept.Value() {
			value = "true"
		}
		return ui.Command{Action: "intercept", Args: []string{value}}
	} else {
		log.Println("Requested intercept status")
		value := "false"
		if intercept.Value() {
			value = "true"
		}
		return ui.Command{Action: "intercept", Args: []string{value}}
	}
}
