package sequence

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "test=foo;")
	}))
	defer ts.Close()

	URL, _ := url.Parse(ts.URL)
	testpayload := `GET / HTTP/1.1
Host: localhost:` + URL.Port() +
		`
User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:54.0) Gecko/20100101 Firefox/54.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
Connection: close
Upgrade-Insecure-Requests: 1


`
	p := newPool(1, []byte(testpayload), false, URL.Host, 1, "test")
	testchan := make(chan string)
	go func() {
		t := time.NewTicker(3 * time.Second)
		select {
		case <-t.C:
			testchan <- "timeout"
		case val := <-p.out:
			testchan <- val
		}
		for _ = range p.out {
		}
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("Stopping pool")
	p.stop()
	cookie := <-testchan
	if cookie == "timeout" {
		t.Error("Test: no cookie received")
	}
	if cookie != "foo" {
		t.Errorf("Got wrong cookie, expected foo but got <%s>", cookie)
	}
	var err error
	select {
	case err = <-p.errors:
	default:
	}
	if err != nil {
		t.Error(err)
	}
}
