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
)

var (
	dec         *json.Decoder
	enc         *json.Encoder
	interceptOn bool
	controls    bool
)

//DOM elements
var (
	jq = jquery.NewJQuery()
	//TODO modify the dom package to provide https://godoc.org/honnef.co/go/js/dom#HTMLTableElement
	//with an insertrow method
	historyTbody      *js.Object
	historyReqBuffer  *DomElement
	historyResBuffer  *DomElement
	proxyBuffer       *DomElement
	endpointIndicator *DomElement
	btn               *DomElement
)

func init() {
	//dom package does not implement table extensions, so we default to bare js
	historyTbody = js.Global.Get("historyTbody")
	historyReqBuffer = GetElementByID("historyReqBuffer")
	historyResBuffer = GetElementByID("historyResBuffer")
	proxyBuffer = GetElementByID("proxybuffer")
	endpointIndicator = GetElementByID("endpointIndicator")
	btn = GetElementByID("interceptToggle")

	//TODO find a way to construct this in the functions declarations
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

//FIXME split in more files
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

	msg.Action = apis.INTERCEPT
	msg.Channel = apis.INTERCEPTSETTINGSCHANNEL

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
		case apis.EDITORCHANNEL:
			proxyBuffer.SetTextContent(string(msg.Payload))
			var text string
			if msg.Args[apis.PAYLOADTYPE] == apis.REQUEST {
				text = "Request for: "
			} else {
				text = "Response from:"
			}
			endpointIndicator.SetTextContent(text + msg.Args[apis.ENDPOINT])
			controls = true

		case apis.INTERCEPTSETTINGSCHANNEL:
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

		case apis.HISTORYCHANNEL:
			switch msg.Action {
			case apis.METADATA:
				var md apis.ReqRespMetaData
				err := json.Unmarshal([]byte(msg.Args[apis.METADATA]), &md)
				if err != nil {
					panic(err)
				}

				if rowMap, ok := tmpHistory[md.ID]; !ok {
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
					tmpHistory[md.ID] = tmp
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

			case apis.FETCH:
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

func onHistoryCellClick() {
	proxyAction(apis.Command{
		Action:  apis.FETCH,
		Channel: apis.HISTORYCHANNEL,
		Args:    map[apis.ArgName]string{apis.ID: js.Global.Get("event").Get("target").Get("parentNode").Get("childNodes").Index(0).Get("textContent").String()},
	}, true)
}
