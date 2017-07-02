// +build js

package main

import (
	"log"

	"github.com/empijei/wapty/ui/apis"
	"github.com/gopherjs/gopherjs/js"
)

func handleEdit(msg apis.Command) {
	proxyBuffer.SetTextContent(string(msg.Payload))
	var text string
	if msg.Args[apis.PAYLOADTYPE] == apis.REQUEST {
		text = "Request for:   "
	} else {
		text = "Response from: "
	}
	endpointIndicator.SetTextContent(text + msg.Args[apis.ENDPOINT])
	controls = true
}

func handleIntercept(msg apis.Command) {
	switch msg.Action {
	case apis.INTERCEPT:
		if msg.Args[apis.ON] == apis.TRUE {
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

func proxyAction(msg apis.Command, ignoreControls bool) {
	log.Printf("Requested action %s", msg.Action)
	if !ignoreControls {
		if !controls {
			return
		}
		controls = false
	}
	log.Printf("Performing action %s", msg.Action)
	err := enc.Encode(msg)
	if err != nil {
		panic(err)
	}
	proxyBuffer.SetTextContent("")
	endpointIndicator.SetTextContent("")
	log.Printf("Action %s invoked", msg.Action)
}

func onForwardOriginalClick() {
	proxyAction(apis.Command{
		Action:  apis.FORWARD,
		Channel: apis.EDITORCHANNEL,
	}, false)
}

func onForwardModifiedClick() {
	proxyAction(apis.Command{
		Action:  apis.EDIT,
		Channel: apis.EDITORCHANNEL,
		Payload: []byte(proxyBuffer.GetTextContent()),
	}, false)
}

func onDropClick() {
	proxyAction(apis.Command{
		Action:  apis.DROP,
		Channel: apis.EDITORCHANNEL,
	}, false)
}

func onProvideResponseClick() {
	proxyAction(apis.Command{
		Action:  apis.PROVIDERESP,
		Channel: apis.EDITORCHANNEL,
		Payload: []byte(proxyBuffer.GetTextContent()),
	}, false)
}

func onToggleInterceptClick() {
	var msg string
	if interceptOn {
		msg = apis.FALSE
	} else {
		msg = apis.TRUE
	}

	var buf string
	if controls && interceptOn {
		buf = proxyBuffer.GetTextContent()
		log.Printf("there is a buffer that will be forwarded, value: %s", buf)
	}

	proxyAction(apis.Command{
		Action:  apis.INTERCEPT,
		Channel: apis.INTERCEPTSETTINGSCHANNEL,
		Args:    map[apis.ArgName]string{apis.ON: msg},
	}, true)

	// If the proxy had a payload when intercept was turned off we assume it was
	// modified
	if buf != "" {
		log.Println("forwarding buffer")
		proxyAction(apis.Command{
			Action:  apis.EDIT,
			Channel: apis.EDITORCHANNEL,
			Payload: []byte(buf),
		}, false)
	}
	interceptOn = !interceptOn
}
