//intercept is meant to handle all the interception of requests and responses,
//including stopping and waiting for edited payloads.
//Every request going through the proxy is parsed and added to the Status by this
//package.
package intercept

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
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
	log.Println("Handling request")
	r := preq.originalRequest
	req, err := httputil.DumpRequest(r, true)
	if err != nil {
		//Something went wrong, abort
		log.Println(err)
		preq.modifiedRequest <- &mayBeRequest{err: err}
		return
	}
	var editedRequest *http.Request
	editedRequestDump, action := editBuffer(REQUEST, req)
	switch action {
	case FORWARD:
		editedRequest = preq.originalRequest
	case EDIT:
		editedRequest, err = http.ReadRequest(bufio.NewReader(bytes.NewReader(editedRequestDump)))
		if err != nil {
			log.Println("Error during edited request parsing, forwarding original request.")
			editedRequest = preq.originalRequest
		}
		//TODO adjust content length
		status.addEditedRequest(preq.id, editedRequestDump)
	case DROP:
		//TODO implement this
		log.Println("Not implemented yet")
		editedRequest = preq.originalRequest
	case PROVIDERESP:
		//TODO implement this
		log.Println("Not implemented yet")
		editedRequest = preq.originalRequest
	default:
		//TODO implement this
		log.Println("Not implemented yet")
		editedRequest = preq.originalRequest
	}

	preq.modifiedRequest <- &mayBeRequest{req: editedRequest}
}

func preProcessRequest(req *http.Request) (autoEdited *http.Request, Id uint, err error) {
	tmp, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println(err.Error())
		return
	}
	Id = newReqResp(tmp)
	//TODO Add autoedit here
	autoEdited = req
	stripHTHHeaders(&(autoEdited.Header))
	//TODO add to status as edited response
	//TODO Return edited one
	//TODO Add auto-resolve hostnames here
	return

	//TODO move this outside
}

func editRequest(req *http.Request, Id uint) (*http.Request, error) {
	//Send request to the dispatchLoop
	ModifiedRequest := make(chan *mayBeRequest)
	RequestQueue <- &pendingRequest{id: Id, originalRequest: req, modifiedRequest: ModifiedRequest}
	log.Println("Request intercepted")
	//Wait for edited request
	mayBeReq := <-ModifiedRequest
	if mayBeReq.err != nil {
		//If edit goes wrong, try to keep going with the original request
		log.Println(mayBeReq.err.Error())
		//FIXME Document this weir behavior or use error properly
		return req, nil
	}
	return mayBeReq.req, nil
}
