package intercept

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"

	"github.com/empijei/Wapty/mitm"
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
	modifiedTransport := responseInterceptor{wrappedRT: http.DefaultTransport}

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
	for {
		select {
		case preq := <-RequestQueue:
			r := preq.originalRequest
			req, err := httputil.DumpRequest(r, true)
			if err != nil {
				//Something went wrong, abort
				preq.modifiedRequest <- &mayBeRequest{err: err}
				break
			}
			_ = ioutil.WriteFile("tmp.request", req, 0644)
			log.Println("Request intercepted, edit it and press enter to continue.")
			_, _ = stdin.ReadString('\n')
			log.Println("Continued")
			//Read the edited request and sent it back to the intercept RequestWrapper
			editedRequestFile, _ := os.Open("tmp.request")                           //TODO chech this error
			editedRequest, _ := http.ReadRequest(bufio.NewReader(editedRequestFile)) //TODO chech this error
			editedRequestDump, _ := httputil.DumpRequest(editedRequest, true)        //TODO chech this error
			status.addEditedRequest(preq.id, &editedRequestDump)
			//Give the request an Id
			editedRequest.Header.Set(idHeader, fmt.Sprintf("%d", preq.id))
			//Mark the request as intercepted
			editedRequest.Header.Set(interceptHeader, "true")
			preq.modifiedRequest <- &mayBeRequest{req: editedRequest}

		case presp := <-ResponseQueue:
			res := presp.originalResponse
			res.ContentLength = -1
			res.Header.Set("Content-Length", "-1")
			rawRes, err := httputil.DumpResponse(res, true)
			status.addResponse(presp.id, &rawRes)
			if err != nil {
				//Something went wrong, abort
				log.Println("Error while dumping response" + err.Error())
				presp.modifiedResponse <- &mayBeResponse{err: err}
				break
			}
			_ = ioutil.WriteFile("tmp.response", rawRes, 0644)
			log.Println("Response intercepted, edit it an press enter to continue.")
			_, _ = stdin.ReadString('\n')
			log.Println("Continued")
			//Read the edited response and sent it back to the responseInterceptor
			editedResponseFile, _ := os.Open("tmp.response") //TODO chech this error
			editedResponseBuffer := bufio.NewReader(editedResponseFile)
			editedResponse, _ := http.ReadResponse(editedResponseBuffer, presp.originalRequest) //TODO chech this error
			editedResponseDump, _ := httputil.DumpResponse(editedResponse, true)                //TODO check this error
			status.addEditedResponse(presp.id, &editedResponseDump)
			//TODO adjust content length Header?
			//fmt.Printf("%s", tmp)
			presp.modifiedResponse <- &mayBeResponse{res: editedResponse, err: err}

		case <-Done:
			return
		}
	}
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
type responseInterceptor struct {
	wrappedRT http.RoundTripper
}

//This is a mock RoundTrip used to intercept responses before they are forwarded by the proxy
func (ri *responseInterceptor) RoundTrip(req *http.Request) (res *http.Response, err error) {
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
