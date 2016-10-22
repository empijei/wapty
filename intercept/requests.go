//intercept is meant to handle all the interception of requests and responses,
//including stopping and waiting for edited payloads.
//Every request going through the proxy is parsed and added to the Status by this
//package.
package intercept

import (
	"bufio"
	"bytes"
	"fmt"
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
		log.Println("Request intercepted")
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
