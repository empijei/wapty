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

type history struct {
	sync.RWMutex
	count    int64
	reqResps map[int64]*ReqResp
}

func parseID(reqId string) (id int64) {
	id, err := strconv.ParseInt(reqId, 10, 64)
	if err != nil {
		panic(err)
	}
	return
}

func (h *history) addEditedRequest(Id int64, rawEditedReq *[]byte) {
	h.RLock()
	h.reqResps[Id].RawEditedReq = rawEditedReq
	h.RUnlock()
}
func (h *history) addResponse(Id int64, rawRes *[]byte) {
	h.RLock()
	h.reqResps[Id].RawRes = rawRes
	h.RUnlock()
}
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

//TODO methods to parse req-resp
type ReqResp struct {
	Id           int64
	RawReq       *[]byte
	RawRes       *[]byte
	RawEditedReq *[]byte
	RawEditedRes *[]byte
}

func newReqResp(rawReq *[]byte) *ReqResp {
	status.Lock()
	status.count += 1
	curReq := status.count
	tmp := &ReqResp{RawReq: rawReq, Id: curReq}
	status.reqResps[curReq] = tmp
	status.Unlock()
	return tmp
}

type mayBeRequest struct {
	req *http.Request
	err error
}
type pendingRequest struct {
	id              int64
	intercepted     bool
	originalRequest *http.Request
	modifiedRequest chan *mayBeRequest
}

type mayBeResponse struct {
	res *http.Response
	err error
}
type pendingResponse struct {
	id               int64
	originalResponse *http.Response
	originalRequest  *http.Request
	modifiedResponse chan *mayBeResponse
}
