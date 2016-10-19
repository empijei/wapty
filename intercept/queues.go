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
var ResponseQueue chan *PendingResponse
var Done chan struct{}
var RequestQueue chan *PendingRequest

var intercept SyncBool

type SyncBool struct {
	sync.Mutex
	value bool
}

func init() {
	RequestQueue = make(chan *PendingRequest)
	ResponseQueue = make(chan *PendingResponse)
	Done = make(chan struct{})
	stdin = bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdin))
	intercept.value = true
}

func MainLoop() {
	ca, err := mitm.LoadCA()
	if err != nil {
		log.Fatal(err)
	}
	go dispatchLoop()
	modifiedTransport := ResponseInterceptor{WrappedRT: http.DefaultTransport}
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
			r := preq.OriginalRequest
			req, err := httputil.DumpRequest(r, true)
			if err != nil {
				preq.ModifiedRequest <- &MayBeRequest{Err: err}
				break
			}
			//fmt.Printf("%s", req)
			_ = ioutil.WriteFile("tmp.request", req, 0644)
			_, _ = stdin.ReadString('\n')
			log.Println("Continued")
			editedRequestFile, _ := os.Open("tmp.request")                           //TODO chech this error
			editedRequest, _ := http.ReadRequest(bufio.NewReader(editedRequestFile)) //TODO chech this error
			editedRequestDump, _ := httputil.DumpRequest(editedRequest, true)        //TODO chech this error
			Status.AddEditedRequest(preq.Id, &editedRequestDump)
			editedRequest.Header.Set(IDHeader, fmt.Sprintf("%d", preq.Id))
			editedRequest.Header.Set(InterceptHeader, "true")
			preq.ModifiedRequest <- &MayBeRequest{Req: editedRequest}

		case presp := <-ResponseQueue:
			res := presp.OriginalResponse
			res.ContentLength = -1
			res.Header.Set("Content-Length", "-1")
			rawRes, err := httputil.DumpResponse(res, true)
			Status.AddResponse(presp.Id, &rawRes)
			if err != nil {
				log.Println("Error while dumping response" + err.Error())
				presp.ModifiedResponse <- &MayBeResponse{Err: err}
				break
			}
			_ = ioutil.WriteFile("tmp.response", rawRes, 0644)
			_, _ = stdin.ReadString('\n')
			log.Println("Continued")
			editedResponseFile, _ := os.Open("tmp.response") //TODO chech this error
			editedResponseBuffer := bufio.NewReader(editedResponseFile)
			editedResponse, _ := http.ReadResponse(editedResponseBuffer, presp.OriginalRequest) //TODO chech this error
			editedResponseDump, _ := httputil.DumpResponse(editedResponse, true)                //TODO check this error
			Status.AddEditedResponse(presp.Id, &editedResponseDump)
			//TODO adjust content length Header?
			//fmt.Printf("%s", tmp)
			presp.ModifiedResponse <- &MayBeResponse{Res: editedResponse, Err: err}

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
		Id := NewReqResp(&tmp).Id
		intercept.Lock()
		intercepted := intercept.value
		intercept.Unlock()
		if !intercepted {
			r.Header.Set(IDHeader, fmt.Sprintf("%d", Id))
			upstream.ServeHTTP(w, r)
			return
		}
		log.Println("Request intercepted")
		//TODO Add autoedit!

		ModifiedRequest := make(chan *MayBeRequest)
		RequestQueue <- &PendingRequest{Id: Id, OriginalRequest: r, ModifiedRequest: ModifiedRequest}
		mayBeReq := <-ModifiedRequest
		if mayBeReq.Err != nil {
			upstream.ServeHTTP(w, r)
			return
		}
		upstream.ServeHTTP(w, mayBeReq.Req)
	})
}

type ResponseInterceptor struct {
	WrappedRT http.RoundTripper
}

func (ri *ResponseInterceptor) RoundTrip(req *http.Request) (res *http.Response, err error) {
	Id := ParseID(req.Header.Get(IDHeader))
	req.Header.Del(IDHeader)
	res, err = ri.WrappedRT.RoundTrip(req)

	//Skip intercept if request was not intercepted
	if req.Header.Get(InterceptHeader) == "" {
		rawRes, dumpErr := httputil.DumpResponse(res, true)
		if dumpErr != nil {
			log.Println(dumpErr.Error())
		} else {
			Status.AddResponse(Id, &rawRes)
		}
		return
	}
	req.Header.Del(InterceptHeader)
	log.Println("Response intercepted")
	if err != nil {
		return
	}
	ModifiedResponse := make(chan *MayBeResponse)
	ResponseQueue <- &PendingResponse{Id: Id, ModifiedResponse: ModifiedResponse, OriginalRequest: req, OriginalResponse: res}
	mayBeRes := <-ModifiedResponse
	res, err = mayBeRes.Res, mayBeRes.Err
	return
}
