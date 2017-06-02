package main

import (
	"encoding/json"
	"fmt"

	"github.com/gopherjs/jquery"
	"github.com/gopherjs/websocket"
)

var jq = jquery.NewJQuery()
var dec *json.Decoder
var enc *json.Encoder

func main() {
	waptyServer, err := websocket.Dial("ws://localhost:8081/ws")
	if err != nil {
		//FIXME handle error
		panic(err)
	}

	fmt.Println("WebSocket connetcted")
	dec = json.NewDecoder(waptyServer)
	enc = json.NewEncoder(waptyServer)

}
