package intercept

import (
	"net/http"
	"strconv"
	"sync"
)

var status history

//Header used to keep track of requests across different routines
const idHeader = "MAPTY-ID"

//Header used to keep track of intercepted requests
const interceptHeader = "MAPTY-Intercept"

func init() {
	status.reqResps = make(map[int64]*ReqResp)
}

//This type is used to represent all the req/resp that went through the proxy
type history struct {
	sync.RWMutex
	count    int64
	reqResps map[int64]*ReqResp
}

//Parses a string into an int64
func parseID(reqId string) (id int64) {
	id, err := strconv.ParseInt(reqId, 10, 64)
	if err != nil {
		panic(err)
	}
	return
}

//Finds the correct Request based on the ID and adds the modified request to it
//This is thread safe
func (h *history) addEditedRequest(Id int64, rawEditedReq *[]byte) {
	h.RLock()
	h.reqResps[Id].RawEditedReq = rawEditedReq
	h.RUnlock()
}

//Finds the correct Request based on the ID and adds the original response to it
//This is thread safe
func (h *history) addResponse(Id int64, rawRes *[]byte) {
	h.RLock()
	h.reqResps[Id].RawRes = rawRes
	h.RUnlock()
}

//Finds the correct Request based on the ID and adds the modified response to it
//This is thread safe
func (h *history) addEditedResponse(Id int64, rawEditedRes *[]byte) {
	h.RLock()
	h.reqResps[Id].RawEditedRes = rawEditedRes
	h.RUnlock()
	//TODO remove this
	//	foo, err := json.MarshalIndent(h.ReqResps[Id], " ", " ")
	//	if err != nil {
	//		log.Println(err.Error())
	//	}
	//	log.Printf("%s", foo)
}

//Represents an item of the proxy history
//TODO methods to parse req-resp
type ReqResp struct {
	//Unique Id in the history
	Id int64
	//Original Request
	RawReq *[]byte
	//Original Response
	RawRes *[]byte
	//Edited Request
	RawEditedReq *[]byte
	//Edited Response
	RawEditedRes *[]byte
}

//Creates a new history item and safely adds it to the status, incrementing the
//current id value
//Returns the id of the newly created item
func newReqResp(rawReq *[]byte) int64 {
	status.Lock()
	status.count += 1
	curReq := status.count
	tmp := &ReqResp{RawReq: rawReq, Id: curReq}
	status.reqResps[curReq] = tmp
	status.Unlock()
	return curReq
}

//represents an *http.Request if err == nil, represents the error otherwise.
type mayBeRequest struct {
	req *http.Request
	err error
}

//a struct used to transmit to the dispatchLoop a requests that waits to be
//edited or forwarded by the user
type pendingRequest struct {
	id              int64
	intercepted     bool
	originalRequest *http.Request
	modifiedRequest chan *mayBeRequest
}

//represents an *http.Response if err == nil, represents the error otherwise.
type mayBeResponse struct {
	res *http.Response
	err error
}

//a struct used to transmit to the dispatchLoop a response that waits to be
//edited or forwarded by the user
type pendingResponse struct {
	id               int64
	originalResponse *http.Response
	originalRequest  *http.Request
	modifiedResponse chan *mayBeResponse
}
