package intercept

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
)

var Status History

const IDHeader = "MAPTY-ID"
const InterceptHeader = "MAPTY-Intercept"

func init() {
	Status.ReqResps = make(map[int64]*ReqResp)
}

type History struct {
	sync.Mutex
	Count    int64
	ReqResps map[int64]*ReqResp
}

func ParseID(reqId string) (Id int64) {
	Id, err := strconv.ParseInt(reqId, 10, 64)
	if err != nil {
		panic(err)
	}
	return
}

func (h *History) AddEditedRequest(Id int64, rawEditedReq *[]byte) {
	h.Lock()
	h.ReqResps[Id].RawEditedReq = rawEditedReq
	h.Unlock()
}
func (h *History) AddResponse(Id int64, rawRes *[]byte) {
	h.Lock()
	h.ReqResps[Id].RawRes = rawRes
	h.Unlock()
}
func (h *History) AddEditedResponse(Id int64, rawEditedRes *[]byte) {
	h.Lock()
	h.ReqResps[Id].RawEditedRes = rawEditedRes
	h.Unlock()
	//TODO remove this
	foo, err := json.MarshalIndent(h.ReqResps[Id], " ", " ")
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("%s", foo)
}

type ReqResp struct {
	Id           int64
	RawReq       *[]byte
	RawRes       *[]byte
	RawEditedReq *[]byte
	RawEditedRes *[]byte
}

func NewReqResp(rawReq *[]byte) *ReqResp {
	Status.Lock()
	Status.Count += 1
	curReq := Status.Count
	tmp := &ReqResp{RawReq: rawReq, Id: curReq}
	Status.ReqResps[curReq] = tmp
	Status.Unlock()
	return tmp
}

//TODO methods to parse req-resp

func main() {
	fmt.Println("vim-go")
}
