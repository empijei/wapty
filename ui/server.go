//Package ui is a general high level representation of all the uis connected to the current
//instance of Wapty. Use this from other packages to read user input and write
//output
package ui

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	rice "github.com/GeertJohan/go.rice"
	"github.com/empijei/wapty/ui/apis"

	"golang.org/x/net/websocket"
)

const BUFSIZE = 1024

var inc = make(chan apis.Command, BUFSIZE)
var outg = make(chan apis.Command, BUFSIZE)

var connMut sync.Mutex
var uiconn io.ReadWriteCloser

func serve(pattern string) {
	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		log.Println("A client has connected")
		defer func() {
			_ = ws.Close()
		}()

		connMut.Lock()
		if uiconn != nil {
			//TODO tell the new UI to GTFO
			log.Println("A UI is already connected")
			return
		}
		uiconn = ws
		connMut.Unlock()

		//This is a blocking call
		handleClient(uiconn)

		connMut.Lock()
		uiconn = nil
		connMut.Unlock()

	}

	http.Handle(pattern, websocket.Handler(onConnected))
	for cmd := range inc {
		subsMutex.RLock()
		for _, out := range subscriptions[cmd.Channel] {
			out.dataCh <- cmd
		}
		subsMutex.RUnlock()
	}
}

func handleClient(uiconn io.ReadWriteCloser) {
	dedicatedchan := make(chan apis.Command)

	//This copyes the commands from the backend to a channel read by the sender goroutine.
	//when the channel is closed it gracefully handles the panic and prevent the last message from being lost.
	go func() {
		var cmd apis.Command
		defer func() {
			if r := recover(); r != nil {
				//cmd was not sent successfully, let's save it
				outg <- cmd
			}
			log.Println("Copyer terminated")
		}()
		for cmd = range outg {
			dedicatedchan <- cmd
		}
	}()

	//Sender goroutine, transmits data from the dedicadedchan to the ui
	go func() {
		enc := json.NewEncoder(uiconn)
		for cmd := range dedicatedchan {
			err := enc.Encode(cmd)
			if err != nil {
				log.Println(err)
				break
			}
		}
		err := uiconn.Close()
		log.Println(err)
		log.Println("Sender terminated")
	}()

	//Takes commands from the ui and sends them to the backend.
	//When the connection is closed this exits and signals the sender to stop.
	dec := json.NewDecoder(uiconn)
	var cmd apis.Command
	for {
		err := dec.Decode(&cmd)
		if err != nil {
			err2 := uiconn.Close()
			log.Println(err2)
			log.Println(err)
			break
		}
		inc <- cmd
	}
	close(dedicatedchan)
	log.Println("A client has disconnected")
}

func send(cmd *apis.Command) {
	outg <- *cmd
}

// MainLoop is the UI's mainloop. It should be run on wapty's start and it will
// not return until an error occours.
func MainLoop() {
	// websocket server
	go serve("/ws")

	// static files
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(rice.MustFindBox("static").HTTPBox())))
	// TODO setup templates

	loadTemplates()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(appPage)
	})
	log.Printf("UI is running on: http://localhost:%d/", 8081)
	log.Fatal(http.ListenAndServe(":8081", nil))

}
