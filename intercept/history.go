package intercept

import (
	"encoding/json"
	"log"
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
}

//This type is used to represent all the req/resp that went through the proxy
type history struct {
	sync.RWMutex
	count    uint
	reqResps []*ReqResp
}

//Parses a string into an uint
func parseID(reqId string) (id uint) {
	sid, err := strconv.Atoi(reqId)
	if err != nil {
		panic(err)
	}
	return uint(sid)
}

//Finds the correct Request based on the ID and adds the modified request to it
//This is thread safe
func (h *history) addEditedRequest(Id uint, rawEditedReq *[]byte) {
	h.RLock()
	h.reqResps[Id].RawEditedReq = rawEditedReq
	h.RUnlock()
}

//Finds the correct Request based on the ID and adds the original response to it
//This is thread safe
func (h *history) addResponse(Id uint, rawRes *[]byte) {
	h.RLock()
	h.reqResps[Id].RawRes = rawRes
	h.RUnlock()
}

//Finds the correct Request based on the ID and adds the modified response to it
//This is thread safe
func (h *history) addEditedResponse(Id uint, rawEditedRes *[]byte) {
	h.RLock()
	h.reqResps[Id].RawEditedRes = rawEditedRes
	h.RUnlock()
	//TODO remove this
	//	foo, err := json.marshalindent(h.reqresps[id], " ", " ")
	//	if err != nil {
	//		log.println(err.error())
	//	}
	//	log.printf("%s", foo)
}

//Dumps the status in the log. This is only meant for debug purposes.
func StatusDump() {
	status.RLock()
	foo, err := json.MarshalIndent(status.reqResps, " ", " ")
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("%s", foo)
	status.RUnlock()
}

//Represents an item of the proxy history
//TODO methods to parse req-resp
type ReqResp struct {
	//Unique Id in the history
	Id uint
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
func newReqResp(rawReq *[]byte) uint {
	status.Lock()
	curReq := status.count
	tmp := &ReqResp{RawReq: rawReq, Id: curReq}
	status.reqResps = append(status.reqResps, tmp)
	status.count += 1
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
	id              uint
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
	id               uint
	originalResponse *http.Response
	originalRequest  *http.Request
	modifiedResponse chan *mayBeResponse
}
