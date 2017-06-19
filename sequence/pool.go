package sequence

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type pool struct {
	wg           sync.WaitGroup
	cookieName   string
	throttledone chan struct{}
	out          chan string
	errors       chan error
	isThrottle   bool

	//fields not to be accessed by workers
	req    http.Request
	ticker *time.Ticker
	done   bool

	//error ratio
	//cookies
}

func newPool(nw int, req http.Request, reqps int, cookieName string) *pool {
	var wg sync.WaitGroup
	wg.Add(nw)

	p := &pool{
		wg:           wg,
		cookieName:   cookieName,
		throttledone: make(chan struct{}, 7*nw),
		out:          make(chan string, 7*nw),
		errors:       make(chan error, 7*nw),
	}

	if reqps > 0 {
		p.isThrottle = true
		p.ticker = time.NewTicker(time.Duration(1000000/reqps) * time.Microsecond)
		go func() {
			//TODO defer recover
			for _ = range p.ticker.C {
				p.throttledone <- struct{}{}
			}
		}()
	}

	for i := 0; i < nw; i++ {
		//TODO clone request
		go doWork(p, req)
	}

	go func() {
		wg.Wait()
		p.done = true
	}()

	return p
}

func doWork(p *pool, req http.Request) {
	defer p.wg.Done()

	//TODO save req body etc

	c := &http.Client{}

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
		resp, err := c.Do(&req)

		//TODO regenerate request
		if err != nil {
			p.errors <- err
		}

		if len(resp.Cookies()) == 0 {
			p.errors <- fmt.Errorf("sequence: no cookie received")
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
