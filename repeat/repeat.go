package repeat

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// DefaultTimeout is the default value for the timeout when creating a new Repeater
var DefaultTimeout = 10 * time.Second

// Item contains the information for a single "Go" of a Repeater
type Item struct {
	Host     string
	TLS      bool
	Request  []byte
	Response []byte
}

// Repeater represents a full history of requests and responses
type Repeater struct {
	m sync.Mutex

	// The Repeater History
	History []Item

	// The timeout to wait before assuming the server will not respond.
	// Default is DefaultTimeout
	Timeout time.Duration
}

// NewRepeater creates a new Repeater with Timeout set to DefaultTimeout
func NewRepeater() *Repeater {
	return &Repeater{
		Timeout: DefaultTimeout,
	}
}

func (r *Repeater) repeat(buf io.Reader, host string, _tls bool) (res io.Reader, err error) {
	r.m.Lock()
	defer r.m.Unlock()
	savedReq := bytes.NewBuffer(nil)
	teebuf := io.TeeReader(buf, savedReq)
	var conn net.Conn
	if _tls {
		//The repeater does not care about certs
		conn, err = tls.Dial("tcp", host, &tls.Config{InsecureSkipVerify: true})
	} else {
		conn, err = net.Dial("tcp", host)
	}
	if err != nil {
		return
	}
	defer func() { _ = conn.Close() }()
	_ = conn.SetDeadline(time.Now().Add(r.Timeout))
	resbuf := bytes.NewBuffer(nil)
	errWrite := make(chan error)

	go func() {
		log.Println("Transmitting the request")
		_, errw := io.Copy(conn, teebuf)
		errWrite <- errw
		log.Println("Request transmitted")
	}()

	log.Println("Reading the response")
	_, err = io.Copy(resbuf, conn)
	log.Println("Response read")
	if tmperr := <-errWrite; tmperr != nil {
		return nil, tmperr
	}
	if err != nil {
		return nil, err
	}
	r.History = append(r.History, Item{Request: savedReq.Bytes(), Response: resbuf.Bytes()})
	return resbuf, nil
}
