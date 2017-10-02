package intercept

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"sync"

	"github.com/empijei/wapty/cli/lg"
	"github.com/empijei/wapty/ui/apis"
)

var status History

func init() {
}

//History is used to represent all the req/resp that went through the proxy
//FIXME!!! implement high-level methods!!!
//FIXME make the fields private and create a dummy object to transmit this
type History struct {
	sync.RWMutex `json:"-"`
	//Remove count, use it only for serialization
	Count    int
	ReqResps []*ReqResp
}

//Finds the correct Request based on the ID and adds the modified request to it
//This is thread safe
func (h *History) addRawEditedRequest(ID int, rawEditedReq []byte) {
	h.RLock()
	h.ReqResps[ID].RawEditedReq = rawEditedReq
	h.RUnlock()
}

//func (h *History) addEditedRequest(ID int, req *http.Request) {
//}

//Finds the correct Request based on the ID and adds the original response to it
//This is thread safe
func (h *History) addRawResponse(ID int, rawRes []byte) {
	h.RLock()
	h.ReqResps[ID].RawRes = rawRes
	h.RUnlock()
}

func (h *History) addResponse(ID int, res *http.Response) {
	tmp, err := httputil.DumpResponse(res, true)
	if err != nil {
		//TODO
		lg.Failure(err.Error())
	}
	h.addRawResponse(ID, tmp)
}

// Finds the correct Request based on the ID and adds the modified response to it
// This is thread safe
func (h *History) addRawEditedResponse(ID int, rawEditedRes []byte) {
	h.RLock()
	h.ReqResps[ID].RawEditedRes = rawEditedRes
	h.RUnlock()
}

//func (h *History) addEditedResponse(ID int, res *http.Response) {}

// StatusDump dumps the status in the log. This is only meant for debug purposes.
func StatusDump(status *History) {
	status.RLock()
	foo, err := json.MarshalIndent(status, " ", " ")
	if err != nil {
		lg.Failure(err.Error())
	}
	lg.Info(foo)
	status.RUnlock()
}

func (h *History) getItem(ID int) *ReqResp {
	h.RLock()
	defer h.RUnlock()
	if ID < h.Count {
		return h.ReqResps[ID]
	}
	return nil
}

//This loop will wait for commands directed to the history control and will
//execute them
func historyLoop() {
	for {
		select {
		case cmd := <-uiHistory.RecChannel():
			switch cmd.Action {
			case apis.HST_DUMP:
				status.RLock()
				dump, err := json.Marshal(status)
				status.RUnlock()
				if err != nil {
					StatusDump(&status)
					panic(err)
				}
				lg.Infof("Dump: %s", dump)
				uiHistory.Send(&apis.Command{Action: "Dump", Payload: dump})
			case apis.HST_FETCH:
				uiHistory.Send(handleFetch(cmd))
			}
		case <-done:
			return
		}
	}
}

func handleFetch(cmd apis.Command) *apis.Command {
	var ID int
	err := cmd.UnpackArgs([]apis.ArgName{apis.ARG_ID}, &ID)
	if err != nil {
		lg.Error(err)
		return apis.Err(err)
	}
	lg.Info("Requested history entry")
	rr := status.getItem(ID)
	buf, err := json.Marshal(rr)
	return &apis.Command{Action: apis.HST_FETCH, Payload: buf}
}

//ReqResp represents an item of the proxy history
//TODO methods to parse req-resp
//TODO create a test that fails if this is different from apis.ReqResp
type ReqResp struct {
	//Unique ID in the history
	ID int
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
	//lg.Info("Locking status for write")
	status.Lock()
	//lg.Info("Locked")
	curReq := status.Count
	tmp := &ReqResp{RawReq: rawReq, ID: curReq, MetaData: &apis.ReqRespMetaData{ID: curReq}}
	status.ReqResps = append(status.ReqResps, tmp)
	status.Count++
	//lg.Info("UnLocking status")
	status.Unlock()
	return curReq
}

func newReqResp(req *http.Request) int {
	tmp, err := httputil.DumpRequest(req, true)
	if err != nil {
		//TODO
		lg.Error(err.Error())
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

// Save saves the status in a json formatted stream
func (h *History) Save(out io.Writer) error {
	h.RLock()
	defer h.RUnlock()

	enc := json.NewEncoder(out)
	err := enc.Encode(h)
	return err
}

// Load loads the status from a json formatted stream
func (h *History) Load(in io.Reader) error {
	var tmp History
	dec := json.NewDecoder(in)
	err := dec.Decode(&tmp)

	if err != nil {
		return err
	}

	h.Lock()
	defer h.Unlock()

	h.ReqResps = tmp.ReqResps
	h.Count = tmp.Count
	return nil
}

// String returns the name of the current package/project
func (h *History) String() string {
	return "Intercept"
}

// GetStatus returns the current status
func GetStatus() *History {
	return &status
}
