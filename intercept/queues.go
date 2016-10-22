//intercept is meant to handle all the interception of requests and responses,
//including stopping and waiting for edited payloads.
//Every request going through the proxy is parsed and added to the Status by this
//package.
package intercept

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"

	"github.com/empijei/Wapty/mitm"
	"github.com/empijei/Wapty/ui"
)

var stdin *bufio.ReadWriter

//Represents the queue of the response to requests that have been intercepted
var ResponseQueue chan *pendingResponse

//Represents the queue of requests that have been intercepted
var RequestQueue chan *pendingRequest

//Not used yet
var Done chan struct{}

//If value is set to true tells the proxy to start the intercept
var intercept SyncBool

type SyncBool struct {
	sync.RWMutex
	value bool
}

func init() {
	RequestQueue = make(chan *pendingRequest)
	ResponseQueue = make(chan *pendingResponse)
	Done = make(chan struct{})
	stdin = bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdin))
	intercept.value = true
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

	//Create the modified transport to intercept responses
	//modifiedTransport := ResponseInterceptor{wrappedRT: http.DefaultTransport} //This uses HTTP2
	noHTTP2Transport := &http.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	modifiedTransport := ResponseInterceptor{wrappedRT: noHTTP2Transport}

	//Creates the mitm.Proxy with the modified transport, the loaded CA and the
	//interceptRequestWrapper
	p := &mitm.Proxy{
		CA: &ca,
		TLSServerConfig: &tls.Config{
			MinVersion: tls.VersionSSL30,
			//CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA},
		},
		Wrap:      interceptRequestWrapper,
		Transport: &modifiedTransport,
	}
	//Starts the mitm.Proxy
	log.Fatal(http.ListenAndServe(":8080", p)) //TODO parametrize this
}

//This loop will keep reading from the RequestQueue and ResponseQueue for new
//intercepted payloads.
//When a request or response is intercepted it is dumped to file to be edited
//and the loop will wait for the user to press enter to continue.
//When a request or response is intercepted and/or modified it is added to the
//History.
func dispatchLoop() {
	uiEditor = ui.Subscribe(EDITORCHANNEL) //FIXME hardcoded string
	for {
		select {
		case preq := <-RequestQueue:
			handleRequest(preq)
		case presp := <-ResponseQueue:
			handleResponse(presp)
		case <-Done:
			return
		}
	}
}

//Called by the dispatchLoop if a request is intercepted
func handleRequest(preq *pendingRequest) {
	r := preq.originalRequest
	req, err := httputil.DumpRequest(r, true)
	if err != nil {
		//Something went wrong, abort
		preq.modifiedRequest <- &mayBeRequest{err: err}
		return
	}
	var editedRequest *http.Request
	editedRequestDump, action := editBuffer(REQUEST, &req)
	switch action {
	case FORWARDED:
		editedRequest = preq.originalRequest
	case EDITED:
		editedRequest, err = http.ReadRequest(bufio.NewReader(bytes.NewReader(*editedRequestDump)))
		if err != nil {
			log.Println("Error during edited request parsing, forwarding original request.")
			editedRequest = preq.originalRequest
		}
		//TODO adjust content length
		status.addEditedRequest(preq.id, editedRequestDump)
	case DROPPED:
		//TODO implement this
		log.Println("Not implemented yet")
		editedRequest = preq.originalRequest
	case RESPPROVIDED:
		//TODO implement this
		log.Println("Not implemented yet")
		editedRequest = preq.originalRequest
	default:
		//TODO implement this
		log.Println("Not implemented yet")
		editedRequest = preq.originalRequest
	}

	//Give the request an Id
	editedRequest.Header.Set(idHeader, fmt.Sprintf("%d", preq.id))
	//Mark the request as intercepted
	editedRequest.Header.Set(interceptHeader, "true")
	preq.modifiedRequest <- &mayBeRequest{req: editedRequest}
}

//Called by the dispatchLoop if a response is intercepted
func handleResponse(presp *pendingResponse) {
	res := presp.originalResponse
	res.ContentLength = -1
	res.Header.Set("Content-Length", "-1")
	rawRes, err := httputil.DumpResponse(res, true)
	status.addResponse(presp.id, &rawRes)
	if err != nil {
		//Something went wrong, abort
		log.Println("Error while dumping response" + err.Error())
		presp.modifiedResponse <- &mayBeResponse{err: err}
		return
	}
	var editedResponse *http.Response
	editedResponseDump, action := editBuffer(RESPONSE, &rawRes)
	switch action {
	case FORWARDED:
		editedResponse = res
	case EDITED:
		editedResponseBuffer := bufio.NewReader(bytes.NewReader(*editedResponseDump))
		editedResponse, err = http.ReadResponse(editedResponseBuffer, presp.originalRequest)
		if err != nil {
			//TODO chech this error and hijack connection to send raw bytes
			log.Println("Error during edited response parsing, forwarding original response.")
			editedResponse = res
		}
		status.addEditedResponse(presp.id, editedResponseDump)
	case DROPPED:
		//TODO implement this
		log.Println("Not implemented yet")
		editedResponse = res
	case RESPPROVIDED:
		//TODO implement this
		log.Println("Action not allowed on Responses")
		editedResponse = res
	default:
		//TODO implement this
		log.Println("Not implemented yet")
		editedResponse = res
	}

	//TODO adjust content length Header?
	//fmt.Printf("%s", tmp)
	presp.modifiedResponse <- &mayBeResponse{res: editedResponse, err: err}

}

//Decorates an http.Handler to make it intercept requests
//see https://www.youtube.com/watch?v=xyDkyFjzFVc
func interceptRequestWrapper(upstream http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmp, err := httputil.DumpRequest(r, true)
		if err != nil {
			upstream.ServeHTTP(w, r)
			log.Println(err.Error())
			return
		}
		Id := newReqResp(&tmp)
		intercept.RLock()
		intercepted := intercept.value
		intercept.RUnlock()
		if !intercepted {
			//If intercept is false just add the request Id and keep going
			r.Header.Set(idHeader, fmt.Sprintf("%d", Id))
			upstream.ServeHTTP(w, r)
			return
		}
		//TODO Add autoedit!

		//Intercept is true, send request to the dispatchLoop
		ModifiedRequest := make(chan *mayBeRequest)
		RequestQueue <- &pendingRequest{id: Id, originalRequest: r, modifiedRequest: ModifiedRequest}
		//Wait for edited request
		mayBeReq := <-ModifiedRequest
		if mayBeReq.err != nil {
			//If edit goes wrong, try to keep going with the original request
			log.Println(mayBeReq.err.Error())
			upstream.ServeHTTP(w, r)
			return
		}
		//Serve the edited request
		upstream.ServeHTTP(w, mayBeReq.req)
	})
}

//This is a struct that respects the net.RoundTripper interface and just wraps
//the original http.RoundTripper
type ResponseInterceptor struct {
	wrappedRT http.RoundTripper
}

//This is a mock RoundTrip used to intercept responses before they are forwarded by the proxy
func (ri *ResponseInterceptor) RoundTrip(req *http.Request) (res *http.Response, err error) {
	//Read request id from header and remove it
	Id := parseID(req.Header.Get(idHeader))
	req.Header.Del(idHeader)
	//Get if the original request was intercepted and remove the header
	intercepted := req.Header.Get(interceptHeader) != ""
	req.Header.Del(interceptHeader)

	//Perform the request
	res, err = ri.wrappedRT.RoundTrip(req)

	if err != nil {
		log.Println("Something went wrong trying to contact the server")
		return
	}

	//Skip intercept if request was not intercepted, only add the response to the Status
	if !intercepted {
		rawRes, dumpErr := httputil.DumpResponse(res, true)
		if dumpErr != nil {
			log.Println(dumpErr.Error())
		} else {
			status.addResponse(Id, &rawRes)
		}
		return
	}

	//Request was intercepted, go throug the intercept/edit process
	ModifiedResponse := make(chan *mayBeResponse)
	ResponseQueue <- &pendingResponse{id: Id, modifiedResponse: ModifiedResponse, originalRequest: req, originalResponse: res}
	mayBeRes := <-ModifiedResponse
	res, err = mayBeRes.res, mayBeRes.err
	return
}
