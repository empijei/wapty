package intercept

import (
	"github.com/empijei/wapty/cli/lg"
	"github.com/empijei/wapty/ui/apis"
)

func settingsLoop() {
	for {
		select {
		case cmd := <-uiSettings.RecChannel():
			lg.Infof("Settings accessed\n")
			switch cmd.Action {
			case apis.STN_INTERCEPT:
				uiSettings.Send(handleIntercept(cmd))
			default:
				//TODO send error?
				lg.Infof("Unknown action: %v\n", cmd.Action)
			}
		case <-done:
			return
		}
	}
}

func handleIntercept(cmd apis.Command) *apis.Command {
	value := apis.ARG_FALSE
	if len(cmd.Args) >= 1 {
		lg.Infof("Requested change intercept status\n")
		intercept.setValue(cmd.Args[apis.ARG_ON] == apis.ARG_TRUE)
		if intercept.value() {
			value = apis.ARG_TRUE
		}
	}
	lg.Infof("Requested intercept status\n")
	if intercept.value() {
		value = apis.ARG_TRUE
	}
	return &apis.Command{
		Action: apis.STN_INTERCEPT,
		Args:   map[apis.ArgName]string{apis.ARG_ON: value},
	}
}
