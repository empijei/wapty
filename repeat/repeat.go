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

var defaultTimeout time.Duration = 10 * time.Second

type RepeatItem struct {
	Host     string
	TLS      bool
	Request  []byte
	Response []byte
}

type Repeater struct {
	m       sync.Mutex
	history []RepeatItem
	timeout time.Duration
}

func NewRepeater() *Repeater {
	return &Repeater{
		timeout: defaultTimeout,
	}
}

func (r *Repeater) Repeat(buf io.Reader, host string, _tls bool) (res io.Reader, err error) {
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
	_ = conn.SetDeadline(time.Now().Add(r.timeout))
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
	r.history = append(r.history, RepeatItem{Request: savedReq.Bytes(), Response: resbuf.Bytes()})
	return resbuf, nil
}
