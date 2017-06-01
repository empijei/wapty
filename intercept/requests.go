//intercept is meant to handle all the interception of requests and responses,
//including stopping and waiting for edited payloads.
//Every request going through the proxy is parsed and added to the Status by this
//package.
package intercept

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/empijei/wapty/ui/apis"
)

//Represents the queue of requests that have been intercepted
var RequestQueue chan *pendingRequest

func init() {
	RequestQueue = make(chan *pendingRequest)
}

//In order for the program to work this should always be started.
//MainLoop is the core of the interceptor. It starts the goroutine that waits
//for new requests and response that have been intercepted and takes action
//based on current configuration.

//Called by the dispatchLoop if a request is intercepted
func handleRequest(preq *pendingRequest) {
	r := preq.originalRequest
	ContentLength := r.ContentLength
	r.ContentLength = -1
	r.Header.Del("Content-Length")
	req, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println("intercept: dumping request " + err.Error())
		preq.modifiedRequest <- &mayBeRequest{err: err}
		return
	}
	var editedRequest *http.Request
	var providedResp *http.Response
	editedRequestDump, action := editBuffer(apis.REQUEST, req, r.URL.Scheme+"://"+r.Host)
	switch action {
	case apis.FORWARD.String():
		r.ContentLength = ContentLength
		editedRequest = r
	case apis.EDIT.String():
		editedRequest, err = editCase(editedRequestDump)
		status.addRawEditedRequest(preq.id, editedRequestDump)
	case apis.DROP.String():
		providedResp = caseDrop()
	case apis.PROVIDERESP.String():
		providedResponseBuffer := bufio.NewReader(bytes.NewReader(editedRequestDump))
		providedResp, err = http.ReadResponse(providedResponseBuffer, preq.originalRequest)
		if err != nil {
			//TODO check this error and hijack connection to send raw bytes
			log.Println("Error during provided response parsingh")
		}
		status.addRawEditedResponse(preq.id, editedRequestDump)
	default:
		//TODO implement this
		log.Println("Not implemented yet")
		editedRequest = preq.originalRequest
	}

	preq.modifiedRequest <- &mayBeRequest{req: editedRequest, res: providedResp, err: err}
}

func preProcessRequest(req *http.Request) (autoEdited *http.Request, Id uint, err error) {
	stripHTHHeaders(&(req.Header))
	Id = newReqResp(req)
	//TODO Add autoedit here
	autoEdited = req
	//FIXME Call this in a "decode" function for requests, like the one used for responses
	//TODO add to status as edited response
	//TODO Return edited one
	//TODO Add auto-resolve hostnames here
	return
}

func editRequest(req *http.Request, Id uint) (*http.Request, *http.Response, error) {
	//Send request to the dispatchLoop
	ModifiedRequest := make(chan *mayBeRequest)
	RequestQueue <- &pendingRequest{id: Id, originalRequest: req, modifiedRequest: ModifiedRequest}
	log.Println("Request intercepted")
	//Wait for edited request
	mayBeReq := <-ModifiedRequest
	if mayBeReq.res != nil {
		return nil, mayBeReq.res, mayBeReq.err
	}
	if mayBeReq.err != nil {
		//If edit goes wrong, try to keep going with the original request
		log.Println(mayBeReq.err.Error())
		//FIXME Document this weir behavior or use error properly
		return req, nil, nil
	}
	return mayBeReq.req, nil, nil
}

func editCase(editedRequestDump []byte) (editedRequest *http.Request, err error) {
	rc := bufio.NewReader(bytes.NewReader(editedRequestDump))
	editedRequest, err = http.ReadRequest(rc)
	if err != nil {
		log.Println("Error during edited request parsing, dunno what to do yet!!!")
		//TODO Default to bare sockets
	}
	//Parsing leftovers, if any, must be the request body
	body, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Println("Error during edited body reading")
		//TODO
	}
	if length := len(body); length != 0 {
		editedRequest.ContentLength = int64(length)
		editedRequest.Body = ioutil.NopCloser(bufio.NewReader(bytes.NewReader(body)))
	}
	return
}
