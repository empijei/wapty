// +build js

package main

//FIXME ui does not respond while receiving big amounts of metadata

import (
	"encoding/json"

	. "github.com/empijei/wapty/ui/apis"
	js "github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
	"github.com/empijei/wapty/cli/lg"
)

var (
	dec         *json.Decoder
	enc         *json.Encoder
	interceptOn bool
)

func send(cmd Command) error {
	return enc.Encode(cmd)
}

func logger(cmd *Command) {
	lg.Infof("Received actions %s on channel %s", cmd.Action, cmd.Channel)
}

func main() {
	document = js.Global.Get("document")

	waptyServer, err := websocket.Dial("ws://localhost:8081/ws")
	if err != nil {
		//FIXME handle error
		panic(err)
	}

	lg.Info("WebSocket connetcted")
	dec = json.NewDecoder(waptyServer)
	enc = json.NewEncoder(waptyServer)

	var msg Command

	msg.Action = STN_INTERCEPT
	msg.Channel = CHN_INTERCEPTSETTINGS

	err = send(msg)
	if err != nil {
		panic(err)
	}

	for {
		var msg Command
		err = dec.Decode(&msg)
		logger(&msg)
		if err != nil {
			panic(err)
		}
		switch msg.Channel {
		case CHN_EDITOR:
			handleEdit(msg)

		case CHN_INTERCEPTSETTINGS:
			handleIntercept(msg)

		case CHN_HISTORY:
			handleHistory(msg)

		case CHN_REPEAT:
			handleRepeat(msg)

		default:
			lg.Error("Unrecognized message")
		}
	}
}
