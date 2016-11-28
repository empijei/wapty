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
	editedRequestDump, action := editBuffer(REQUEST, req)
	switch action {
	case FORWARD:
		r.ContentLength = ContentLength
		editedRequest = r
	case EDIT:
		editedRequest, err = editCase(editedRequestDump)
		status.addRawEditedRequest(preq.id, editedRequestDump)
	case DROP:
		providedResp = caseDrop()
	case PROVIDERESP:
		//TODO implement this
		log.Println("Not implemented yet")
		editedRequest = r
	default:
		//TODO implement this
		log.Println("Not implemented yet")
		editedRequest = preq.originalRequest
	}

	preq.modifiedRequest <- &mayBeRequest{req: editedRequest, res: providedResp, err: err}
}

func preProcessRequest(req *http.Request) (autoEdited *http.Request, Id uint, err error) {
	Id = newReqResp(req)
	//TODO Add autoedit here
	autoEdited = req
	//FIXME Call this in a "decode" function for requests, like the one used for responses
	//move this to beginning of function
	stripHTHHeaders(&(autoEdited.Header))
	//TODO add to status as edited response
	//TODO Return edited one
	//TODO Add auto-resolve hostnames here
	return

	//TODO move this outside
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
	//Adjust Content-Length
	//Whiile Responses default to having a body,
	//Requests will not unles Transfer-Encoding or Content-Length is specified
	//http://greenbytes.de/tech/webdav/draft-ietf-httpbis-p1-messaging-26.html#message.body
	//First attempt, did not work
	//FIXME not all browsers post removing \r, detect the case!
	//TODO fix the last 2 lines? \n\r?

	//tmpSplit := bytes.SplitN(editedRequestDump, []byte("\n"), 2)
	//if len(tmpSplit) > 1 {
	//if tmpSplit[1][0] != byte('\r') {
	//tmpSplit = bytes.SplitN(editedRequestDump, []byte("\n\n"), 2)
	////Fixing \n\r for headers
	//tmpSplit[0] = bytes.Replace(tmpSplit[0], []byte("\n"), []byte("\n\r"), -1)
	//} else {
	//tmpSplit = bytes.SplitN(editedRequestDump, []byte("\n\r\n\r"), 2)
	//}
	//}

	//if len(tmpSplit) > 1 {
	//cl := strconv.Itoa(len(tmpSplit[1]))
	//log.Println("Content Length:" + cl)
	//var tmp [][]byte
	//tmp = append(tmp, tmpSplit[0], []byte("\n\rContent-Length: "+cl+"\n\r\n\r"), tmpSplit[1])
	//editedDump := bytes.Join(tmp, nil)
	//editedRequest, err = http.ReadRequest(bufio.newreader(bytes.newreader(editeddump)))
	//} else {
	//editedRequest, err = http.ReadRequest(bufio.NewReader(bytes.NewReader(editedRequestDump)))
	//}

	//Original, does not work
	//editedRequest, err = http.ReadRequest(bufio.NewReader(bytes.NewReader(editedRequestDump)))

	rc := bufio.NewReader(bytes.NewReader(editedRequestDump))
	editedRequest, err = http.ReadRequest(rc)
	if err != nil {
		log.Println("Error during edited request parsing, dunno wat todo!!!")
		//TODO Default to bare sockets
	}
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
