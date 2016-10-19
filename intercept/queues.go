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
	sync.Mutex
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
	ca, err := mitm.LoadCA()
	if err != nil {
		log.Fatal(err)
	}
	go dispatchLoop()
	modifiedTransport := responseInterceptor{wrappedRT: http.DefaultTransport}
	p := &mitm.Proxy{
		CA: &ca,
		TLSServerConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			//CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA},
		},
		Wrap:      interceptRequestWrapper,
		Transport: &modifiedTransport,
	}
	log.Fatal(http.ListenAndServe(":8080", p))
}

func dispatchLoop() {
	for {
		select {
		case preq := <-RequestQueue:
			r := preq.originalRequest
			req, err := httputil.DumpRequest(r, true)
			if err != nil {
				preq.modifiedRequest <- &mayBeRequest{err: err}
				break
			}
			//fmt.Printf("%s", req)
			_ = ioutil.WriteFile("tmp.request", req, 0644)
			_, _ = stdin.ReadString('\n')
			log.Println("Continued")
			editedRequestFile, _ := os.Open("tmp.request")                           //TODO chech this error
			editedRequest, _ := http.ReadRequest(bufio.NewReader(editedRequestFile)) //TODO chech this error
			editedRequestDump, _ := httputil.DumpRequest(editedRequest, true)        //TODO chech this error
			status.addEditedRequest(preq.id, &editedRequestDump)
			editedRequest.Header.Set(idHeader, fmt.Sprintf("%d", preq.id))
			editedRequest.Header.Set(interceptHeader, "true")
			preq.modifiedRequest <- &mayBeRequest{req: editedRequest}

		case presp := <-ResponseQueue:
			res := presp.originalResponse
			res.ContentLength = -1
			res.Header.Set("Content-Length", "-1")
			rawRes, err := httputil.DumpResponse(res, true)
			status.addResponse(presp.id, &rawRes)
			if err != nil {
				log.Println("Error while dumping response" + err.Error())
				presp.modifiedResponse <- &mayBeResponse{err: err}
				break
			}
			_ = ioutil.WriteFile("tmp.response", rawRes, 0644)
			_, _ = stdin.ReadString('\n')
			log.Println("Continued")
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

func interceptRequestWrapper(upstream http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmp, err := httputil.DumpRequest(r, true)
		if err != nil {
			upstream.ServeHTTP(w, r)
			log.Println(err.Error())
			return
		}
		Id := newReqResp(&tmp).Id
		intercept.Lock()
		intercepted := intercept.value
		intercept.Unlock()
		if !intercepted {
			r.Header.Set(idHeader, fmt.Sprintf("%d", Id))
			upstream.ServeHTTP(w, r)
			return
		}
		log.Println("Request intercepted")
		//TODO Add autoedit!

		ModifiedRequest := make(chan *mayBeRequest)
		RequestQueue <- &pendingRequest{id: Id, originalRequest: r, modifiedRequest: ModifiedRequest}
		mayBeReq := <-ModifiedRequest
		if mayBeReq.err != nil {
			upstream.ServeHTTP(w, r)
			return
		}
		upstream.ServeHTTP(w, mayBeReq.req)
	})
}

type responseInterceptor struct {
	wrappedRT http.RoundTripper
}

//This is a mock RoundTrip used to intercept responses before they are forwarded
//by the proxy
func (ri *responseInterceptor) RoundTrip(req *http.Request) (res *http.Response, err error) {
	Id := parseID(req.Header.Get(idHeader))
	req.Header.Del(idHeader)
	res, err = ri.wrappedRT.RoundTrip(req)

	//Skip intercept if request was not intercepted
	if req.Header.Get(interceptHeader) == "" {
		rawRes, dumpErr := httputil.DumpResponse(res, true)
		if dumpErr != nil {
			log.Println(dumpErr.Error())
		} else {
			status.addResponse(Id, &rawRes)
		}
		return
	}
	req.Header.Del(interceptHeader)
	log.Println("Response intercepted")
	if err != nil {
		return
	}
	ModifiedResponse := make(chan *mayBeResponse)
	ResponseQueue <- &pendingResponse{id: Id, modifiedResponse: ModifiedResponse, originalRequest: req, originalResponse: res}
	mayBeRes := <-ModifiedResponse
	res, err = mayBeRes.res, mayBeRes.err
	return
}
