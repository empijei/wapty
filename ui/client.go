package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/empijei/wapty/ui/apis"

	"golang.org/x/net/websocket"
)

//TODO this is not good practice
const channelBufSize = 100

var maxID int

// client represents the server-side representation of a connected UI
type client struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan *apis.Command
	doneCh chan bool
}

// newClient instantiates a new Client
func newClient(ws *websocket.Conn, server *Server) *client {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxID++
	ch := make(chan *apis.Command, channelBufSize)
	doneCh := make(chan bool)

	return &client{maxID, ws, server, ch, doneCh}
}

// Conn returns the Client underlying websocket
func (c *client) Conn() *websocket.Conn {
	return c.ws
}

func (c *client) write(msg *apis.Command) {
	select {
	case c.ch <- msg:
	default:
		c.server.DelClient(c)
		err := fmt.Errorf("client %d is disconnected", c.id)
		c.server.Err(err)
	}
}

// done signals the Client to stop
func (c *client) done() {
	close(c.doneCh)
}

func (c *client) listen() {
	go c.listenWrite()
	c.listenRead()
}

func (c *client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			err := websocket.JSON.Send(c.ws, msg)
			if err != nil {
				log.Println("Error while sending message to websocket client.")
				close(c.doneCh)
			}

		case <-c.doneCh:
			c.server.DelClient(c)
			return
		}
	}
}

// Listen read request via chanel
func (c *client) listenRead() {
	log.Println("Listening read from client")
	dec := json.NewDecoder(c.ws)
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.DelClient(c)
			return

		// read data from websocket connection
		default:
			var msg apis.Command
			//HOLY SMOKES THIS DOES NOT SUPPORT MULTIPLE FRAMES
			//Are you serious? https://github.com/golang/go/issues/7632
			//Unmarshal accepts a []byte instead of an io.Reader() so
			//that prevents an easy fix for this problem without refactoring
			//and use a json.NewDecoder(io.Reader) instead of Unmarshal([]byte,...)
			//
			//err := websocket.JSON.Receive(c.ws, &msg)
			err := dec.Decode(&msg)
			if err == io.EOF {
				close(c.doneCh)
			} else if err != nil {
				c.server.Err(err)
				log.Println(msg)
				close(c.doneCh)
			} else {
				//TODO check if action != nil
				//log.Println("Received ", msg)
				c.server.msgReceived(&msg)
			}
		}
	}
}
