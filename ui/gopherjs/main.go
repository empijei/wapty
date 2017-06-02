package main

import (
	"encoding/json"
	"log"

	"github.com/empijei/wapty/ui/apis"
	"github.com/gopherjs/jquery"
	"github.com/gopherjs/websocket"

	"honnef.co/go/js/dom"
)

var jq = jquery.NewJQuery()
var dec *json.Decoder
var enc *json.Encoder
var document dom.Document

func main() {
	document = dom.GetWindow().Document()

	waptyServer, err := websocket.Dial("ws://localhost:8081/ws")
	if err != nil {
		//FIXME handle error
		panic(err)
	}

	log.Println("WebSocket connetcted")
	dec = json.NewDecoder(waptyServer)
	enc = json.NewEncoder(waptyServer)

	var msg apis.Command

	msg.Action = apis.INTERCEPT.String()
	msg.Channel = apis.SETTINGSCHANNEL

	err = enc.Encode(msg)
	if err != nil {
		panic(err)
	}

	//TODO to improve
	var historyTbody = document.GetElementByID("historyTbody")
	var historyReqBuffer = document.GetElementByID("historyReqBuffer")
	var historyResBuffer = document.GetElementByID("historyResBuffer")
	var proxyBuffer = document.GetElementByID("proxybuffer")
	var endpointIndicator = document.GetElementByID("endpointIndicator")
	var btn = document.GetElementByID("interceptToggle")

	var interceptOn bool
	var controls bool

	for {
		err = dec.Decode(&msg)
		if err != nil {
			panic(err)
		}
		switch msg.Channel {
		case apis.EDITORCHANNEL.String():
			proxyBuffer.SetNodeValue(string(msg.Payload))
			var text string
			if msg.Args[0] == apis.REQUEST.String() {
				text = "Request for: "
			} else {
				text = "Response for:"
			}
			endpointIndicator.SetTextContent(text + msg.Args[1])
			controls = true

		case apis.SETTINGSCHANNEL.String():
			switch msg.Action {
			case apis.INTERCEPT.String():
				if msg.Args[0] == true {
					btn.Class().Remove("btn-danger")
					btn.Class().Set([]string{"btn-success"})
					btn.SetTextContent("Intercept is on")
					interceptOn = true
				} else {
					btn.Class().Remove("btn-success")
					btn.Class().Set([]string{"btn-danger"})
					btn.SetTextContent("Intercept is off")
					interceptOn = false
				}
			}

		case apis.HISTORYCHANNEL.String():
			switch msg.Action {
			case apis.METADATA.String():
			case apis.FETCH.String():
			}
		default:
			log.Println("Unrecognized message")
		}
	}
}
