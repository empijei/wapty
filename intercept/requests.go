package intercept

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/empijei/wapty/cli/lg"
	"github.com/empijei/wapty/ui/apis"
)

// RequestQueue represents the queue of requests that have been intercepted
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
		lg.Errorf("intercept: dumping request %s\n", err.Error())
		preq.modifiedRequest <- &mayBeRequest{err: err}
		return
	}
	var editedRequest *http.Request
	var providedResp *http.Response
	editedRequestDump, action := editBuffer(apis.PLD_REQUEST, req, r.URL.Scheme+"://"+r.Host)
	switch action {
	case apis.EDT_FORWARD:
		r.ContentLength = ContentLength
		editedRequest = r
	case apis.EDT_EDIT:
		editedRequest, err = editCase(editedRequestDump)
		status.addRawEditedRequest(preq.id, editedRequestDump)
	case apis.EDT_DROP:
		providedResp = caseDrop()
	case apis.EDT_PROVIDERESP:
		providedResponseBuffer := bufio.NewReader(bytes.NewReader(editedRequestDump))
		providedResp, err = http.ReadResponse(providedResponseBuffer, preq.originalRequest)
		if err != nil {
			//TODO check this error and hijack connection to send raw bytes
			lg.Errorf("Error during provided response parsingh\n")
		}
		status.addRawEditedResponse(preq.id, editedRequestDump)
	default:
		//TODO implement this
		lg.Errorf("Not implemented yet\n")
		editedRequest = preq.originalRequest
	}

	preq.modifiedRequest <- &mayBeRequest{req: editedRequest, res: providedResp, err: err}
}

func preProcessRequest(req *http.Request) (autoEdited *http.Request, ID int, err error) {
	stripHTHHeaders(&(req.Header))
	ID = newReqResp(req)
	//TODO Add autoedit here
	autoEdited = req
	//FIXME Call this in a "decode" function for requests, like the one used for responses
	//TODO add to status as edited response
	//TODO Return edited one
	//TODO Add auto-resolve hostnames here
	return
}

func editRequest(req *http.Request, ID int) (*http.Request, *http.Response, error) {
	//Send request to the dispatchLoop
	ModifiedRequest := make(chan *mayBeRequest)
	RequestQueue <- &pendingRequest{id: ID, originalRequest: req, modifiedRequest: ModifiedRequest}
	lg.Infof("Request intercepted\n")
	//Wait for edited request
	mayBeReq := <-ModifiedRequest
	if mayBeReq.res != nil {
		return nil, mayBeReq.res, mayBeReq.err
	}
	if mayBeReq.err != nil {
		//If edit goes wrong, try to keep going with the original request
		lg.Errorf("%s\n", mayBeReq.err.Error())
		//FIXME Document this weir behavior or use error properly
		return req, nil, nil
	}
	return mayBeReq.req, nil, nil
}

func editCase(editedRequestDump []byte) (editedRequest *http.Request, err error) {
	rc := bufio.NewReader(bytes.NewReader(editedRequestDump))
	editedRequest, err = http.ReadRequest(rc)
	if err != nil {
		lg.Errorf("Error during edited request parsing, dunno what to do yet!!!\n")
		//TODO Default to bare sockets
	}
	//Parsing leftovers, if any, must be the request body
	body, err := ioutil.ReadAll(rc)
	if err != nil {
		lg.Errorf("Error during edited body reading\n")
		//TODO
	}
	if length := len(body); length != 0 {
		editedRequest.ContentLength = int64(length)
		editedRequest.Body = ioutil.NopCloser(bufio.NewReader(bytes.NewReader(body)))
	}
	return
}
