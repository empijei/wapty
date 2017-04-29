package repeat

import (
	"bytes"
	"io"
	"net"
)

type RepeatItem struct {
	Request  []byte
	Response []byte
}

type Repeater struct {
	history []RepeatItem
}

func (r *Repeater) Repeat(buf io.Reader, host string) (res io.Reader, err error) {
	savedReq := bytes.NewBuffer(nil)
	teebuf := io.TeeReader(buf, savedReq)
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return
	}
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
