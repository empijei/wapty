package intercept

import (
	"github.com/empijei/wapty/cli/lg"
	"github.com/empijei/wapty/ui/apis"
)

func settingsLoop() {
	for {
		select {
		case cmd := <-uiSettings.RecChannel():
			lg.Debug("Settings accessed")
			switch cmd.Action {
			case apis.STN_INTERCEPT:
				uiSettings.Send(handleIntercept(cmd))
			default:
				//TODO send error?
				lg.Error("Unknown action: %v", cmd.Action)
			}
		case <-done:
			return
		}
	}
}

func handleIntercept(cmd apis.Command) *apis.Command {
	value := apis.ARG_FALSE
	if len(cmd.Args) >= 1 {
		lg.Debug("Requested change intercept status")
		intercept.setValue(cmd.Args[apis.ARG_ON] == apis.ARG_TRUE)
		if intercept.value() {
			value = apis.ARG_TRUE
		}
	}
	lg.Debug("Requested intercept status")
	if intercept.value() {
		value = apis.ARG_TRUE
	}
	return &apis.Command{
		Action: apis.STN_INTERCEPT,
		Args:   map[apis.ArgName]string{apis.ARG_ON: value},
	}
}
