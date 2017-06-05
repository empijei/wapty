package intercept

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"sync"

	"github.com/empijei/wapty/ui/apis"
)

var status History

func init() {
}

//FIXME!!! implement high-level methods!!!
//This type is used to represent all the req/resp that went through the proxy
//FIXME make the fields private and create a dummy object to transmit this
type History struct {
	sync.RWMutex
	//Remove count, use it only for serialization
	Count    int
	ReqResps []*ReqResp
}

//Finds the correct Request based on the ID and adds the modified request to it
//This is thread safe
func (h *History) addRawEditedRequest(Id int, rawEditedReq []byte) {
	h.RLock()
	h.ReqResps[Id].RawEditedReq = rawEditedReq
	h.RUnlock()
}

//func (h *History) addEditedRequest(Id int, req *http.Request) {
//}

//Finds the correct Request based on the ID and adds the original response to it
//This is thread safe
func (h *History) addRawResponse(Id int, rawRes []byte) {
	h.RLock()
	h.ReqResps[Id].RawRes = rawRes
	h.RUnlock()
}

func (h *History) addResponse(Id int, res *http.Response) {
	tmp, err := httputil.DumpResponse(res, true)
	if err != nil {
		//TODO
		log.Println(err.Error())
	}
	h.addRawResponse(Id, tmp)
}

//Finds the correct Request based on the ID and adds the modified response to it
//This is thread safe
func (h *History) addRawEditedResponse(Id int, rawEditedRes []byte) {
	h.RLock()
	h.ReqResps[Id].RawEditedRes = rawEditedRes
	h.RUnlock()
}

//func (h *History) addEditedResponse(Id int, res *http.Response) {}

//Dumps the status in the log. This is only meant for debug purposes.
func StatusDump(status History) {
	status.RLock()
	foo, err := json.MarshalIndent(status, " ", " ")
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("%s", foo)
	status.RUnlock()
}

func (h *History) getItem(Id int) *ReqResp {
	h.RLock()
	defer h.RUnlock()
	if Id < h.Count {
		return h.ReqResps[Id]
	} else {
		return nil
	}
}

//This loop will wait for commands directed to the history control and will
//execute them
func historyLoop() {
	for {
		select {
		case cmd := <-uiHistory.DataChannel:
			switch cmd.Action {
			case apis.DUMP.String():
				status.RLock()
				dump, err := json.Marshal(status)
				status.RUnlock()
				if err != nil {
					StatusDump(status)
					panic(err)
				}
				log.Printf("Dump: %s\n", dump)
				uiHistory.Send(apis.Command{Action: "Dump", Payload: dump})
			case apis.FETCH.String():
				uiHistory.Send(handleFetch(cmd))
			}
		case <-done:
			return
		}
	}
}

func handleFetch(cmd apis.Command) apis.Command {
	if len(cmd.Args) >= 1 {
		log.Println("Requested history entry")
		Id, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return apis.Command{Action: "Error", Args: []string{"Invalid argument to FETCH"}}
		}
		rr := status.getItem(Id)
		buf, err := json.Marshal(rr)
		return apis.Command{Action: apis.FETCH.String(), Payload: buf}
	} else {
		log.Println("Missing argument for FETCH")
		return apis.Command{Action: "Error", Args: []string{"Missing argument for FETCH"}}
	}
}

//Represents an item of the proxy history
//TODO methods to parse req-resp
//TODO create a test that fails if this is different from apis.ReqResp
type ReqResp struct {
	//Unique Id in the history
	Id int
	//Meta Data about both Req and Resp
	MetaData *apis.ReqRespMetaData
	//Original Request
	RawReq []byte
	//Original Response
	RawRes []byte
	//Edited Request
	RawEditedReq []byte
	//Edited Response
	RawEditedRes []byte
}

//Creates a new history item and safely adds it to the status, incrementing the
//current id value
//Returns the id of the newly created item
func newRawReqResp(rawReq []byte) int {
	//log.Println("Locking status for write")
	status.Lock()
	//log.Println("Locked")
	curReq := status.Count
	tmp := &ReqResp{RawReq: rawReq, Id: curReq, MetaData: apis.NewReqRespMetaData(curReq)}
	status.ReqResps = append(status.ReqResps, tmp)
	status.Count += 1
	//log.Println("UnLocking status")
	status.Unlock()
	return curReq
}

func newReqResp(req *http.Request) int {
	tmp, err := httputil.DumpRequest(req, true)
	if err != nil {
		//TODO
		log.Println(err.Error())
	}
	return newRawReqResp(tmp)
}

//represents an *http.Request if err == nil, represents the error otherwise.
type mayBeRequest struct {
	req *http.Request
	res *http.Response
	err error
}

//a struct used to transmit to the dispatchLoop a requests that waits to be
//edited or forwarded by the user
type pendingRequest struct {
	id              int
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
	id               int
	originalResponse *http.Response
	originalRequest  *http.Request
	modifiedResponse chan *mayBeResponse
}
