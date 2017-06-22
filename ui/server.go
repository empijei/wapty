package ui

import (
	"log"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/empijei/wapty/ui/apis"

	"golang.org/x/net/websocket"
)

// Ui server.
type Server struct {
	pattern   string
	clients   map[int]*client
	addCh     chan *client
	delCh     chan *client
	sendAllCh chan *apis.Command
	doneCh    chan bool
	errCh     chan error
}

// Create new ui server.
func NewServer(pattern string) *Server {
	clients := make(map[int]*client)
	addCh := make(chan *client)
	delCh := make(chan *client)
	sendAllCh := make(chan *apis.Command)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}

func (s *Server) AddClient(c *client) {
	s.addCh <- c
}

func (s *Server) DelClient(c *client) {
	s.delCh <- c
}

func (s *Server) SendAllClients(msg *apis.Command) {
	s.sendAllCh <- msg
}

func (s *Server) Done() {
	close(s.doneCh)
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendAllClients(msg *apis.Command) {
	for _, c := range s.clients {
		c.write(msg)
	}
}

func (s *Server) msgReceived(msg *apis.Command) {
	Receive(*msg)
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := newClient(ws, s)
		s.AddClient(client)
		client.listen()
	}

	//TODO only listen for localhost
	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")
			//TODO do something, send status?

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)

		// broadcast message for all clients
		case msg := <-s.sendAllCh:
			//log.Println("Send all:", msg)
			s.sendAllClients(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}

func writeLoop(s *Server) {
	oChan := Connect()
	for msg := range oChan.Channel() {
		s.SendAllClients(&msg)
	}
}

// MainLoop is the UI's mainloop. It should be run on wapty's start and it will
// not return until an error occours.
func MainLoop() {
	// websocket server
	server := NewServer("/ws")
	go server.Listen()

	go writeLoop(server)

	// static files
	//http.Handle("/", http.FileServer(assetFS()))
	http.Handle("/", http.FileServer(rice.MustFindBox("webroot").HTTPBox()))

	//This is a dirty workaround for the websocket package not reensembling frames
	//http.HandleFunc("/edit", func(rw http.ResponseWriter, req *http.Request) {
	//decoder := json.NewDecoder(req.Body)
	//defer func() { _ = req.Body.Close() }()
	//var cmd apis.Command
	//err := decoder.Decode(&cmd)
	//if err != nil {
	//log.Println(err)
	//}
	//server.msgReceived(&cmd)
	//})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
