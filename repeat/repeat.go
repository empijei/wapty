package repeat

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"sync"
	"time"
)

type RepeatItem struct {
	Host     string
	TLS      bool
	Request  []byte
	Response []byte
}

type Repeater struct {
	m       sync.Mutex
	history []RepeatItem
}

func NewRepeater() *Repeater {
	return &Repeater{}
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
	//TODO check this works
	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	resbuf := bytes.NewBuffer(nil)
	errWrite := make(chan error)

	go func() {
		_, errw := io.Copy(conn, teebuf)
		errWrite <- errw
	}()

	_, err = io.Copy(resbuf, conn)
	if tmperr := <-errWrite; tmperr != nil {
		return nil, tmperr
	}
	//TODO add savedreq to repeater history
	r.history = append(r.history, RepeatItem{Request: savedReq.Bytes(), Response: resbuf.Bytes()})
	return
}
