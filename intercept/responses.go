package intercept

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/empijei/wapty/ui/apis"
)

//ResponseQueue represents the queue of the response to requests that have been intercepted
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
	if err != nil {
		log.Println("intercept: dumping response " + err.Error())
		presp.modifiedResponse <- &mayBeResponse{err: err}
		return
	}
	var editedResponse *http.Response
	editedResponseDump, action := editBuffer(apis.RESPONSE, rawRes, presp.originalRequest.URL.Scheme+"://"+presp.originalRequest.Host)
	switch action {
	case apis.FORWARD:
		res.ContentLength = ContentLength
		res.Header.Set("Content-Length", strconv.Itoa(int(ContentLength)))
		editedResponse = res
	case apis.EDIT, apis.PROVIDERESP:
		editedResponseBuffer := bufio.NewReader(bytes.NewReader(editedResponseDump))
		editedResponse, err = http.ReadResponse(editedResponseBuffer, presp.originalRequest)
		if err != nil {
			//TODO check this error and hijack connection to send raw bytes
			log.Println("Error during edited response parsing, forwarding original response.")
			res.ContentLength = ContentLength
			editedResponse = res
		}
		status.addRawEditedResponse(presp.id, editedResponseDump)
	case apis.DROP:
		editedResponse = caseDrop()
	default:
		//TODO implement this
		log.Println("Not implemented yet")
	}
	presp.modifiedResponse <- &mayBeResponse{res: editedResponse, err: err}
}
func preProcessResponse(req *http.Request, res *http.Response, ID int) *http.Response {
	res = decodeResponse(res)
	//Skip intercept if request was not intercepted, just add the response to the Status
	status.addResponse(ID, res)
	//TODO autoEdit here
	//TODO add to status as edited if autoedited
	return res
}
func editResponse(req *http.Request, res *http.Response, ID int) (*http.Response, error) {
	//Request was intercepted, go through the intercept/edit process
	//TODO use the autoedited one to edit
	ModifiedResponse := make(chan *mayBeResponse)
	ResponseQueue <- &pendingResponse{id: ID, modifiedResponse: ModifiedResponse, originalRequest: req, originalResponse: res}
	mayBeRes := <-ModifiedResponse
	return mayBeRes.res, mayBeRes.err
}

//Making use of the net.http package to remove all the encoding by exausting the
//request body and replacing it with a io.ReadCloser with the complete response.
//This takes care of Transfer-Encoding and Content-Encoding
func decodeResponse(res *http.Response) *http.Response {
	defer func() { _ = res.Body.Close() }()
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return res
	}
	res.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	res.TransferEncoding = nil
	stripHTHHeaders(&(res.Header))
	res.ContentLength = int64(len(buf))
	return res
}

func caseDrop() (res *http.Response) {
	return GenerateResponse("Interceptor", "Response was dropped", 418)
}
