package main

//FIXME ui does not respond while receiving big amounts of metadata

import (
	"encoding/json"
	"log"

	"github.com/empijei/wapty/ui/apis"
	"github.com/gopherjs/websocket"
)

var (
	dec         *json.Decoder
	enc         *json.Encoder
	interceptOn bool
)

func send(cmd apis.Command) error {
	return enc.Encode(cmd)
}

func logger(cmd *apis.Command) {
	log.Printf("Received actions %s on channel %s", cmd.Action, cmd.Channel)
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

	msg.Action = apis.INTERCEPT
	msg.Channel = apis.INTERCEPTSETTINGSCHANNEL

	err = send(msg)
	if err != nil {
		panic(err)
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
			handleEdit(msg)

		case apis.INTERCEPTSETTINGSCHANNEL:
			handleIntercept(msg)

		case apis.HISTORYCHANNEL:
			handleHistory(msg)

		case apis.REPEATCHANNEL:
			handleRepeat(msg)

		default:
			log.Println("Unrecognized message")
		}
	}
}
