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

//Represents the queue of the response to requests that have been intercepted
var ResponseQueue chan *pendingResponse

func init() {
	ResponseQueue = make(chan *pendingResponse)
}

//Called by the dispatchLoop if a response is intercepted
func handleResponse(presp *pendingResponse) {
	res := presp.originalResponse
	//res = uncompress(res)
	ContentLength := res.ContentLength
	res.ContentLength = -1
	res.Header.Del("Content-Length")
	rawRes, err := httputil.DumpResponse(res, true)
	status.addResponse(presp.id, rawRes)
	if err != nil {
		//Something went wrong, abort
		log.Println("Error while dumping response" + err.Error())
		presp.modifiedResponse <- &mayBeResponse{err: err}
		return
	}
	var editedResponse *http.Response
	editedResponseDump, action := editBuffer(RESPONSE, rawRes)
	switch action {
	case FORWARD:
		res.ContentLength = ContentLength
		editedResponse = res
	case EDIT:
		editedResponseBuffer := bufio.NewReader(bytes.NewReader(editedResponseDump))
		editedResponse, err = http.ReadResponse(editedResponseBuffer, presp.originalRequest)
		if err != nil {
			//TODO check this error and hijack connection to send raw bytes
			log.Println("Error during edited response parsing, forwarding original response.")
			editedResponse = res
		}
		status.addEditedResponse(presp.id, editedResponseDump)
	case DROP:
		//TODO implement this
		log.Println("Not implemented yet")
		editedResponse = res
	case PROVIDERESP:
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

func editResponse(req *http.Request, res *http.Response, intercepted bool, Id uint) (*http.Response, error) {
	res = decode(res)
	//Skip intercept if request was not intercepted, only add the response to the Status
	rawRes, dumpErr := httputil.DumpResponse(res, true)
	if dumpErr != nil {
		log.Println(dumpErr.Error())
	} else {
		status.addResponse(Id, rawRes)
	}
	//TODO autoEdit here
	//TODO add to status as edited if autoedited
	if !intercepted {
		return res, dumpErr
	}

	//Request was intercepted, go through the intercept/edit process
	//TODO use the autoedited one to edit
	ModifiedResponse := make(chan *mayBeResponse)
	ResponseQueue <- &pendingResponse{id: Id, modifiedResponse: ModifiedResponse, originalRequest: req, originalResponse: res}
	mayBeRes := <-ModifiedResponse
	return mayBeRes.res, mayBeRes.err
}

//Making use of the net.http package to remove all the encoding by exausting the
//request body and replacing it with a io.ReadCloser with the complete response.
//This takes care of Transfer-Encoding and Content-Encoding
func decode(res *http.Response) *http.Response {
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return res
	}
	res.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	res.TransferEncoding = nil
	res.Header.Del("Content-Encoding")
	res.ContentLength = int64(len(buf))
	return res
}
