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

//Represents the queue of the response to requests that have been intercepted
var ResponseQueue chan *pendingResponse

func init() {
	ResponseQueue = make(chan *pendingResponse)
}

//Called by the dispatchLoop if a response is intercepted
func handleResponse(presp *pendingResponse) {
	res := presp.originalResponse
	ContentLength := res.ContentLength
	res.ContentLength = -1
	res.Header.Del("Content-Length")
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
		res.ContentLength = ContentLength
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

	StatusDump()
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
