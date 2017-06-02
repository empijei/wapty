package main

//FIXME ui does not respond while receiving big amounts of metadata

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/empijei/wapty/ui/apis"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/gopherjs/websocket"

	"honnef.co/go/js/dom"
)

var (
	dec         *json.Decoder
	enc         *json.Encoder
	interceptOn bool
	controls    bool
)

//DOM elements
var (
	jq       = jquery.NewJQuery()
	document dom.Document
	//TODO modify the dom package to provide https://godoc.org/honnef.co/go/js/dom#HTMLTableElement
	//with an insertrow method
	historyTbody      *js.Object
	historyReqBuffer  dom.Element
	historyResBuffer  dom.Element
	proxyBuffer       dom.Element
	endpointIndicator dom.Element
	btn               dom.Element
)

func init() {
	document = dom.GetWindow().Document()
	//dom package does not implement table extensions, so we default to bare js
	historyTbody = js.Global.Get("historyTbody")
	historyReqBuffer = document.GetElementByID("historyReqBuffer")
	historyResBuffer = document.GetElementByID("historyResBuffer")
	proxyBuffer = document.GetElementByID("proxybuffer")
	endpointIndicator = document.GetElementByID("endpointIndicator")
	btn = document.GetElementByID("interceptToggle")

	js.Global.Set("proxy", map[string]interface{}{
		"onDropClick":            onDropClick,
		"onForwardModifiedClick": onForwardModifiedClick,
		"onForwardOriginalClick": onForwardOriginalClick,
		"onProvideResponseClick": onProvideResponseClick,
		"onToggleInterceptClick": onToggleInterceptClick,
		"onHistoryCellClick":     onHistoryCellClick,
	})

	hth := js.Global.Get("historyHeader")

	//This is used to make the ui adapt to backend changes in metadata
	val := reflect.Indirect(reflect.ValueOf(apis.ReqRespMetaData{}))
	for i := 0; i < val.Type().NumField(); i++ {
		hth.Call("insertCell", -1).Set("innerText", val.Type().Field(i).Name)
	}
}

func main() {
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
	msg.Channel = apis.SETTINGSCHANNEL.String()

	err = enc.Encode(msg)
	if err != nil {
		panic(err)
	}

	tmpHistory := make(map[int]map[string]*js.Object)
	logger := func(cmd *apis.Command) {
		log.Printf("Received actions %s on channel %s", cmd.Action, cmd.Channel)
	}
	for {
		var msg apis.Command
		err = dec.Decode(&msg)
		logger(&msg)
		if err != nil {
			panic(err)
		}
		switch msg.Channel {
		case apis.EDITORCHANNEL.String():
			proxyBuffer.SetTextContent(string(msg.Payload))
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
				if msg.Args[0] == "true" {
					btn.Class().Set([]string{"btn", "btn-success"})
					btn.SetTextContent("Intercept is on")
					interceptOn = true
				} else {
					btn.Class().Set([]string{"btn", "btn-danger"})
					btn.SetTextContent("Intercept is off")
					interceptOn = false
				}
			}

		case apis.HISTORYCHANNEL.String():
			switch msg.Action {
			case apis.METADATA.String():
				var md apis.ReqRespMetaData
				err := json.Unmarshal([]byte(msg.Args[0]), &md)
				if err != nil {
					panic(err)
				}

				if rowMap, ok := tmpHistory[md.Id]; !ok {
					row := historyTbody.Call("insertRow", -1)
					tmp := make(map[string]*js.Object)

					val := reflect.Indirect(reflect.ValueOf(md))
					for i := 0; i < val.Type().NumField(); i++ {
						typeField := val.Type().Field(i)
						cell := row.Call("insertCell", -1)
						valueField := val.Field(i).Interface()
						cell.Set("innerText", fmt.Sprintf("%v", valueField))
						tmp[typeField.Name] = cell
					}
					tmpHistory[md.Id] = tmp
				} else {
					val := reflect.Indirect(reflect.ValueOf(md))
					for i := 0; i < val.Type().NumField(); i++ {
						typeField := val.Type().Field(i)
						cell := rowMap[typeField.Name]
						valueField := val.Field(i).Interface()
						cell.Set("innerText", fmt.Sprintf("%v", valueField))
					}
					//FIXME this is commented because it looks like the page
					//receives the same metadata multiple times.
					//delete(tmpHistory, md.Id)
				}

			case apis.FETCH.String():
				var rr apis.ReqResp
				err := json.Unmarshal(msg.Payload, &rr)
				if err != nil {
					panic(err)
				}
				//TODO check if is printable, otherwise show hex
				historyReqBuffer.SetTextContent(string(rr.RawReq))
				historyResBuffer.SetTextContent(string(rr.RawRes))
			}
		default:
			log.Println("Unrecognized message")
		}
	}
}

func proxyAction(msg apis.Command, ignoreControls bool) {
	log.Printf("Invoking action %s", msg.Action)
	if !ignoreControls {
		if !controls {
			return
		}
		controls = false
	}
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
		Action:  apis.FORWARD.String(),
		Channel: apis.EDITORCHANNEL.String(),
	}, false)
}

func onForwardModifiedClick() {
	proxyAction(apis.Command{
		Action:  apis.EDIT.String(),
		Channel: apis.EDITORCHANNEL.String(),
		Payload: []byte(proxyBuffer.NodeValue()),
	}, false)
}

func onDropClick() {
	proxyAction(apis.Command{
		Action:  apis.DROP.String(),
		Channel: apis.EDITORCHANNEL.String(),
	}, false)
}

func onProvideResponseClick() {
	proxyAction(apis.Command{
		Action:  apis.PROVIDERESP.String(),
		Channel: apis.EDITORCHANNEL.String(),
		Payload: []byte(proxyBuffer.NodeValue()),
	}, false)
}

func onToggleInterceptClick() {
	var msg string
	if interceptOn {
		msg = "false"
	} else {
		msg = "true"
	}

	proxyAction(apis.Command{
		Action:  apis.INTERCEPT.String(),
		Channel: apis.SETTINGSCHANNEL.String(),
		Args:    []string{msg},
	}, true)

	interceptOn = !interceptOn
}

func onHistoryCellClick() {
	proxyAction(apis.Command{
		Action:  apis.FETCH.String(),
		Channel: apis.HISTORYCHANNEL.String(),
		Args: []string{dom.WrapEvent(
			js.Global.Get("event")).Target().ParentNode().ChildNodes()[0].TextContent(),
		},
	}, true)
}
