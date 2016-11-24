//intercept is meant to handle all the interception of requests and responses,
//including stopping and waiting for edited payloads.
//Every request going through the proxy is parsed and added to the Status by this
//package.
package intercept

import (
	"crypto/tls"
	"log"
	"net/http"
	"sync"

	"github.com/empijei/Wapty/mitm"
	"github.com/empijei/Wapty/ui"
)

//Not used yet
var done chan struct{}

//If value is set to true tells the proxy to start the intercept
var intercept SyncBool

var uiSettings *ui.Subscription

type SyncBool struct {
	sync.RWMutex
	value bool
}

func init() {
	done = make(chan struct{})
	//intercept.value = true
	uiSettings = ui.Subscribe(SETTINGSCHANNEL)
}

//In order for the program to work this should always be started.
//MainLoop is the core of the interceptor. It starts the goroutine that waits
//for new requests and response that have been intercepted and takes action
//based on current configuration.
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
	noHTTP2Transport := &http.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	//noHTTP2Transport.DisableCompression = false
	modifiedTransport := Interceptor{wrappedRT: noHTTP2Transport}

	//Creates the mitm.Proxy with the modified transport, the loaded CA and the
	//interceptRequestWrapper
	p := &mitm.Proxy{
		CA: &ca,
		TLSServerConfig: &tls.Config{
			MinVersion: tls.VersionSSL30,
			//CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA},
		},
		//FIXME disabled for debug purposes
		//Wrap:      interceptRequestWrapper,
		Transport: &modifiedTransport,
	}
	//Starts the mitm.Proxy
	log.Println(http.ListenAndServe(":8080", p)) //TODO parametrize this
	done <- struct{}{}
	done <- struct{}{}
	done <- struct{}{}
}

//This loop will keep reading from the RequestQueue and ResponseQueue for new
//intercepted payloads.
//When a request or response is intercepted it is dumped to file to be edited
//and the loop will wait for the user to press enter to continue.
//When a request or response is intercepted and/or modified it is added to the
//History.
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

//This is a struct that respects the net.RoundTripper interface and just wraps
//the original http.RoundTripper
type Interceptor struct {
	wrappedRT http.RoundTripper
}
