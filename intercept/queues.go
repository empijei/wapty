//Package intercept is meant to handle all the interception of requests and responses,
//including stopping and waiting for edited payloads.
//Every request going through the proxy is parsed and added to the Status by this
//package.
package intercept

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/empijei/wapty/mitm"
	"github.com/empijei/wapty/ui"
	"github.com/empijei/wapty/ui/apis"
)

//Not used yet
var done chan struct{}

//If value is set to true tells the proxy to start the intercept
var intercept syncBool

var uiSettings ui.Subscription

func init() {
	done = make(chan struct{})
	//intercept.value = true
	uiSettings = ui.Subscribe(apis.CHN_INTERCEPTSETTINGS)
}

// syncBool is just used as a thread safe bool. As of 19/06/2017 the sync/atomic
// package does not provide boolean operations
type syncBool struct {
	sync.RWMutex
	val bool
}

func (s *syncBool) value() bool {
	s.RLock()
	defer s.RUnlock()
	return s.val
}

func (s *syncBool) setValue(v bool) {
	s.Lock()
	s.val = v
	s.Unlock()
}

// MainLoop is the core of the interceptor.
// In order for the normal lifecycle program to work this should always be started.
// It starts the goroutine that waits for new requests and response that have
// been intercepted and takes action based on current configuration.
func MainLoop() {
	//Load Certificate authority
	ca, err := mitm.LoadCA()
	if err != nil {
		log.Fatal(err)
	}

	//Call dispatchloop on other goroutine
	go dispatchLoop()

	//Run History interactions
	go historyLoop()

	//Listen for settings changes
	go settingsLoop()

	//Create the modified transport to intercept responses
	//modifiedTransport := ResponseInterceptor{wrappedRT: http.DefaultTransport} //This uses HTTP2
	wrappedTransport := &http.Transport{
		TLSNextProto:        make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		TLSHandshakeTimeout: 5 * time.Second, //TODO make this a variable
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second, //TODO make this a variable
		}).Dial,
	}
	//noHTTP2Transport.DisableCompression = false
	modifiedTransport := Interceptor{wrappedRT: wrappedTransport}

	//Creates the mitm.Proxy with the modified transport and the loaded CA
	p := &mitm.Proxy{
		CA: &ca,
		TLSServerConfig: &tls.Config{
			MinVersion: tls.VersionSSL30,
		},
		//Wrap:      interceptRequestWrapper,
		Transport: &modifiedTransport,
	}

	log.Printf("Proxy is running on localhost:%d", 8080)
	//Starts the mitm.Proxy
	log.Println(http.ListenAndServe(":8080", p)) //TODO parametrize this and allow for closure
	close(done)
}

//This loop will keep reading from the RequestQueue and ResponseQueue for new
//intercepted payloads.
func dispatchLoop() {
	for {
		select {
		case preq := <-RequestQueue:
			handleRequest(preq)
		case presp := <-ResponseQueue:
			handleResponse(presp)
		case <-done:
			return
		}
	}
}
