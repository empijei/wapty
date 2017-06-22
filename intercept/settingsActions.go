package intercept

import (
	"log"

	"github.com/empijei/wapty/ui/apis"
)

func settingsLoop() {
	for {
		select {
		case cmd := <-uiSettings.RecChannel():
			log.Println("Settings accessed")
			switch cmd.Action {
			case apis.INTERCEPT:
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

func handleIntercept(cmd apis.Command) apis.Command {
	value := apis.FALSE
	if len(cmd.Args) >= 1 {
		log.Println("Requested change intercept status")
		intercept.setValue(cmd.Args[apis.ON] == apis.TRUE)
		if intercept.value() {
			value = apis.TRUE
		}
	}
	log.Println("Requested intercept status")
	if intercept.value() {
		value = apis.TRUE
	}
	return apis.Command{
		Action: apis.INTERCEPT,
		Args:   map[apis.Param]string{apis.ON: value},
	}
}
