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

var maxId int

type Client struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan *apis.Command
	doneCh chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxId++
	ch := make(chan *apis.Command, channelBufSize)
	doneCh := make(chan bool)

	return &Client{maxId, ws, server, ch, doneCh}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(msg *apis.Command) {
	select {
	case c.ch <- msg:
	default:
		c.server.DelClient(c)
		err := fmt.Errorf("client %d is disconnected", c.id)
		c.server.Err(err)
	}
}

func (c *Client) Done() {
	close(c.doneCh)
}

func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

func (c *Client) listenWrite() {
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
func (c *Client) listenRead() {
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
