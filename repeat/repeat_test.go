package repeat

import (
	"bytes"
	"io/ioutil"
	"net"
	"testing"
	"time"
)

type RepTest struct {
	in  []byte
	out []byte
}

var RepeatTests = []RepTest{
	{
		[]byte(`GET /success.txt HTTP/1.1
Host: detectportal.firefox.com
User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:53.0) Gecko/20100101 Firefox/53.0
Accept: */*
Accept-Language: en-US,en;q=0.5
Cache-Control: no-cache
Pragma: no-cache
Connection: close

	`),
		[]byte(`HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 8
Last-Modified: Mon, 15 May 2017 18:04:40 GMT
ETag: "ae780585f49b94ce1444eb7d28906123"
Accept-Ranges: bytes
Server: AmazonS3
X-Amz-Cf-Id: MnfbeXeS3ep60gjgpK6jEZF5WYcQix8AeNXFZBLf8RpVEOC1kWBUUQ==
Cache-Control: no-cache, no-store, must-revalidate
Date: Tue, 23 May 2017 09:29:16 GMT
Connection: close

success
	`),
	},
}

func listener(t *testing.T, testChan chan RepTest, input chan []byte, l net.Listener) {
	defer func() { _ = l.Close() }()
	for c, err := l.Accept(); err == nil; c, err = l.Accept() {
		t.Log("Got incoming connection")
		tt := <-testChan
		buf := make([]byte, len(tt.in))
		n, err := c.Read(buf)
		if err != nil {
			t.Error(err)
		}
		t.Logf("Read from connection %d bytes", n)
		for tmp := 0; n < len(tt.in) && err == nil; {
			tmp, err = c.Read(buf[n+tmp:])
			if tmp != 0 {
				t.Logf("Read from connection %d more bytes", tmp)
			}
			n += tmp
		}
		if err != nil {
			t.Error(err)
		}
		if n != len(tt.in) {
			t.Errorf("Expected read of %d but actually read %d bytes.", len(tt.in), n)
		}
		input <- buf
		n, err = c.Write(tt.out)
		if err != nil {
			t.Error(err)
		}
		if n != len(tt.out) {
			t.Errorf("Should have written %d bytes but wrote %d", len(tt.out), n)
		}
		_ = c.Close()
	}
}

func TestRepeatPlain(t *testing.T) {
	testChan := make(chan RepTest, 2)
	input := make(chan []byte, 2)
	t.Log("Listening on port 12321")
	l, err := net.Listen("tcp", ":12321")
	if err != nil {
		t.Fatal(err)
	}
	go listener(t, testChan, input, l)
	for _, tt := range RepeatTests {
		testChan <- tt
		DefaultTimeout = 1 * time.Second
		r := NewRepeater()
		buf := bytes.NewBuffer(tt.in)
		res, _, err := r.repeat(buf, "localhost:12321", false)
		if err != nil {
			t.Error(err)
			return
		}
		resBuf, err := ioutil.ReadAll(res)
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(resBuf, tt.out) != 0 {
			t.Errorf("Expected <%s> but got <%s>", string(tt.out), string(resBuf))
		}
		in := <-input
		if bytes.Compare(in, tt.in) != 0 {
			t.Errorf("Expected <%s> but got <%s>", string(tt.in), string(in))
		}
	}
}
