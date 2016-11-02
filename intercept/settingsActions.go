package intercept

import (
	"log"

	"github.com/empijei/Wapty/ui"
)

//String used to recognize commands directed to this module
const SETTINGSCHANNEL = "proxy/intercept/options"

func settingsLoop() {
	for {
		cmd := uiSettings.Read()
		log.Println("Settings change requested")
		switch cmd.Action {
		case "intercept":
			if len(cmd.Args) >= 1 {
				log.Println("Requested change intercept status")
				intercept.Lock()
				intercept.value = cmd.Args[0] == "true"
				intercept.Unlock()
			} else {
				log.Println("Requested intercept status")
				intercept.RLock()
				value := "false"
				if intercept.value == true {
					value = "true"
				}
				intercept.RUnlock()
				ui.Send(ui.Command{Channel: SETTINGSCHANNEL, Action: "intercept", Args: []string{value}})
			}
		}
	}
}
