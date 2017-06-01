// +build js
package main

import "github.com/gopherjs/websocket"

func main() {
	_, err := websocket.Dial("ws://localhost:8081/ws")
	if err != nil {
		//FIXME handle error
		panic(err)
	}
}
