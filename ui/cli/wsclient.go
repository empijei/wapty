package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/empijei/Wapty/intercept"
	"github.com/empijei/Wapty/ui"
	"golang.org/x/net/websocket"
)

var ws *websocket.Conn
var serverChannel chan ui.Command
var stdin *bufio.ReadWriter

func init() {
	serverChannel = make(chan ui.Command)
	stdin = bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdin))
}

func main() {
	go wsLoop()
	cli()
}

func cli() {
	for cmd := range serverChannel {
		_ = ioutil.WriteFile("tmp.swp", *cmd.Payload, 0644)
		log.Println("Payload intercepted, edit it and press enter to continue.")
		_, _ = stdin.ReadString('\n')
		log.Println("Continued")
		payload, _ := ioutil.ReadFile("tmp.swp") //TODO chech this error
		args := ui.Args(map[string]string{"action": intercept.EDITED.String()})
		err := websocket.JSON.Send(ws, ui.Command{Args: args, Channel: intercept.EDITORCHANNEL, Payload: &payload})
		if err != nil {
			panic(err)
		}
	}
}

func wsLoop() {
	var url = "ws://localhost:8081/ws"
	var origin = "http://localhost/"
	var err error
	ws, err = websocket.Dial(url, "", origin)
	if err != nil {
		panic(err)
	}
	for {
		var msg ui.Command
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			panic(err)
		}
		serverChannel <- msg
	}

}
