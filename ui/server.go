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
			return
		}
		uiconn = ws
		connMut.Unlock()
		handleClient()
		log.Println("A client has disconnected")
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

func handleClient() {
	go func() {
		enc := json.NewEncoder(uiconn)
		for cmd := range outg {
			err := enc.Encode(cmd)
			if err != nil {
				connMut.Lock()
				defer connMut.Unlock()
				if uiconn != nil {
					_ = uiconn.Close()
					uiconn = nil
				}
				log.Println(err)
				return
			}
		}
	}()
	dec := json.NewDecoder(uiconn)
	var cmd apis.Command
	for {
		err := dec.Decode(&cmd)
		if err != nil {
			connMut.Lock()
			defer connMut.Unlock()
			if uiconn != nil {
				_ = uiconn.Close()
				uiconn = nil
			}
			log.Println(err)
			return
		}
		inc <- cmd
	}
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
		w.Write(appPage)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
