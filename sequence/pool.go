package sequence

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/empijei/cli/lg"
)

var timeout time.Duration = 3 * time.Second

type pool struct {
	wg           sync.WaitGroup
	cookieName   string
	throttledone chan struct{}
	out          chan string
	errors       chan error
	isThrottle   bool
	req          []byte
	_tls         bool
	host         string

	//fields not to be accessed by workers
	ticker *time.Ticker
	done   bool

	//error ratio
	//cookies
}

//TODO rethrottle
//TODO expose info on status and errors
//TODO add counter, close when done

func newPool(nw int, req []byte, _tls bool, host string, reqps int, cookieName string) *pool {
	var wg sync.WaitGroup
	wg.Add(nw)

	p := &pool{
		wg:           wg,
		cookieName:   cookieName,
		throttledone: make(chan struct{}, 7*nw),
		out:          make(chan string, 7*nw),
		errors:       make(chan error, 7*nw),
		req:          req,
		_tls:         _tls,
		host:         host,
	}
	if reqps > 0 {
		p.isThrottle = true
		p.ticker = time.NewTicker(time.Duration(1000000/reqps) * time.Microsecond)
		go func() {
			// TODO fix this, this is not good practice
			defer func() {
				_ = recover()
			}()
			for range p.ticker.C {
				p.throttledone <- struct{}{}
			}
		}()
	}

	for i := 0; i < nw; i++ {
		go doWork(p)
	}

	go func() {
		wg.Wait()
		close(p.out)
		p.done = true
	}()

	return p
}

func (p *pool) stop() {
	close(p.throttledone)
}

func doWork(p *pool) {
	defer p.wg.Done()

	for {
		//throttledone is used only as a "done" channel if worker is unthrottled
		//and is used to acquire execution tokens if the worker is throttled
		if p.isThrottle {
			select {
			case _, ok := <-p.throttledone:
				if !ok {
					return
				}
			}
		} else {
			select {
			case _, ok := <-p.throttledone:
				if !ok {
					return
				}
			default:
			}
		}

		//perform req
		resp, err := doReq(p.req, p._tls, p.host)

		//TODO regenerate request
		if err != nil {
			p.errors <- err
			continue
		}

		if len(resp.Cookies()) == 0 {
			p.errors <- fmt.Errorf("sequence: no cookie received")
			continue
		}

		//get the right cookie
		found := false
		for _, cookie := range resp.Cookies() {
			if cookie.Name == p.cookieName {
				//send cookie over channel
				p.out <- cookie.Value
				found = true
				break
			}
		}
		if !found {
			p.errors <- fmt.Errorf("sequence: no session cookie received")
		}
	}
}

func doReq(buf []byte, _tls bool, host string) (resp *http.Response, err error) {
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
	_ = conn.SetDeadline(time.Now().Add(timeout))
	resbuf := bytes.NewBuffer(nil)
	errWrite := make(chan error)

	teebuf := bytes.NewBuffer(buf)

	go func() {
		_, errw := io.Copy(conn, teebuf)
		errWrite <- errw
	}()

	_, err = io.Copy(resbuf, conn)
	if tmperr := <-errWrite; tmperr != nil {
		err = tmperr
		return
	}
	if err != nil {
		return
	}

	lg.Debug(string(resbuf.Bytes()))
	return http.ReadResponse(bufio.NewReader(resbuf), nil)
}
