package intercept

import (
	"log"

	"github.com/empijei/Wapty/ui"
)

//String used to recognize commands directed to this module
const SETTINGSCHANNEL = "proxy/intercept/options"

func settingsLoop() {
	for {
		select {
		case cmd := <-uiSettings.DataChannel:
			log.Println("Settings accessed")
			switch cmd.Action {
			case "intercept":
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
		intercept.Lock()
		intercept.value = cmd.Args[0] == "true" || cmd.Args[0] == "on"
		value := "false"
		if intercept.value {
			value = "true"
		}
		intercept.Unlock()
		return ui.Command{Action: "intercept", Args: []string{value}}
	} else {
		log.Println("Requested intercept status")
		intercept.RLock()
		value := "false"
		if intercept.value {
			value = "true"
		}
		intercept.RUnlock()
		return ui.Command{Action: "intercept", Args: []string{value}}
	}
}
