// +build js

package main

import (
	"log"

	. "github.com/empijei/wapty/ui/apis"
	"github.com/gopherjs/gopherjs/js"
	"github.com/empijei/wapty/cli/lg"
)

func handleEdit(msg Command) {
	proxyBuffer.SetTextContent(string(msg.Payload))
	var text string
	if msg.Args[ARG_PAYLOADTYPE] == PLD_REQUEST {
		text = "Request for:   "
	} else {
		text = "Response from: "
	}
	endpointIndicator.SetTextContent(text + msg.Args[ARG_ENDPOINT])
	controls = true
}

func handleIntercept(msg Command) {
	switch msg.Action {
	case STN_INTERCEPT:
		if msg.Args[ARG_ON] == ARG_TRUE {
			btn.ToggleClass("btn-danger", "btn-success")
			btn.SetTextContent("Intercept is on")
			interceptOn = true
		} else {
			btn.ToggleClass("btn-success", "btn-danger")
			btn.SetTextContent("Intercept is off")
			interceptOn = false
		}
	}
}

//DOM Stuff
var (
	proxyBuffer       *DomElement
	endpointIndicator *DomElement
	btn               *DomElement
	controls          bool
)

func init() {
	proxyBuffer = GetElementByID("proxybuffer")
	endpointIndicator = GetElementByID("endpointIndicator")
	btn = GetElementByID("interceptToggle")
	js.Global.Set("proxy", map[string]interface{}{
		"onDropClick":            onDropClick,
		"onForwardModifiedClick": onForwardModifiedClick,
		"onForwardOriginalClick": onForwardOriginalClick,
		"onProvideResponseClick": onProvideResponseClick,
		"onToggleInterceptClick": onToggleInterceptClick,
		"onHistoryCellClick":     onHistoryCellClick,
	})
}

func proxyAction(msg Command, ignoreControls bool) {
	lg.Infof("Requested action %s", msg.Action)
	if !ignoreControls {
		if !controls {
			return
		}
		controls = false
	}
	lg.Infof("Performing action %s", msg.Action)
	err := enc.Encode(msg)
	if err != nil {
		panic(err)
	}
	proxyBuffer.SetTextContent("")
	endpointIndicator.SetTextContent("")
	lg.Infof("Action %s invoked", msg.Action)
}

func onForwardOriginalClick() {
	proxyAction(Command{
		Action:  EDT_FORWARD,
		Channel: CHN_EDITOR,
	}, false)
}

func onForwardModifiedClick() {
	proxyAction(Command{
		Action:  EDT_EDIT,
		Channel: CHN_EDITOR,
		Payload: []byte(proxyBuffer.GetTextContent()),
	}, false)
}

func onDropClick() {
	proxyAction(Command{
		Action:  EDT_DROP,
		Channel: CHN_EDITOR,
	}, false)
}

func onProvideResponseClick() {
	proxyAction(Command{
		Action:  EDT_PROVIDERESP,
		Channel: CHN_EDITOR,
		Payload: []byte(proxyBuffer.GetTextContent()),
	}, false)
}

func onToggleInterceptClick() {
	var msg string
	if interceptOn {
		msg = ARG_FALSE
	} else {
		msg = ARG_TRUE
	}

	var buf string
	if controls && interceptOn {
		buf = proxyBuffer.GetTextContent()
		lg.Infof("there is a buffer that will be forwarded, value: %s", buf)
	}

	proxyAction(Command{
		Action:  STN_INTERCEPT,
		Channel: CHN_INTERCEPTSETTINGS,
		Args:    map[ArgName]string{ARG_ON: msg},
	}, true)

	// If the proxy had a payload when intercept was turned off we assume it was
	// modified
	if buf != "" {
		lg.Infof("forwarding buffer\n")
		proxyAction(Command{
			Action:  EDT_EDIT,
			Channel: CHN_EDITOR,
			Payload: []byte(buf),
		}, false)
	}
	interceptOn = !interceptOn
}
