package intercept

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/empijei/Wapty/ui"
)

type ReqRespMetaData struct {
	Id          uint
	Host        string
	Method      string
	Path        string
	Params      bool
	Edited      bool
	Status      string
	Length      int64
	ContentType string
	Extension   string
	TLS         bool
	IP          string
	Port        string
	Cookies     string
	Time        string
	/*
		Port!
		Title (maybe not?)
		Comment (user-defined)
	*/
}

func newReqRespMetaData(Id uint) *ReqRespMetaData {
	return &ReqRespMetaData{Id: Id}
}

//DISCLAIMER use original req AFTER editing the new one
//And use it from a thread that has a readlock on the status
func (rr *ReqResp) parseRequest(req *http.Request) {
	this := rr.MetaData
	this.Host = req.Host
	this.Method = req.Method
	this.Path = req.URL.Path
	if len(req.Form) == 0 {
		_ = req.ParseForm()
	}
	this.Params = len(req.Form) > 0
	//this supposes to alread have a RLock on the status.
	this.Edited = status.ReqResps[this.Id].RawEditedReq != nil
	tmp := strings.Split(this.Path, ".")
	if !strings.Contains(tmp[len(tmp)-1], "/") {
		this.Extension = tmp[len(tmp)-1]
	}
	ips, err := net.LookupHost(strings.Split(this.Host, ":")[0])
	if err == nil && len(ips) >= 1 {
		this.IP = ips[0]
		if len(ips) >= 2 {
			this.Port = ips[1]
		} else {
			log.Println("Port not specified")
		}
	}
	this.Time = time.Now().String()
	sendMetaData(this)
}

//DISCLAIMER use original res AFTER editing the new one
//And use it from a thread that has a readlock on the status
func (rr *ReqResp) parseResponse(res *http.Response) {
	this := rr.MetaData
	if !this.Edited {
		//this supposes to alread have a RLock on the status.
		this.Edited = status.ReqResps[this.Id].RawEditedRes != nil
	}
	this.Status = res.Status
	this.Length = res.ContentLength
	this.ContentType = res.Header.Get("Content-Type")
	this.TLS = res.TLS != nil
	tmp := res.Cookies()
	for _, cookie := range tmp {
		this.Cookies += cookie.String() + "; "
	}
	sendMetaData(this)
}

func sendMetaData(metaData *ReqRespMetaData) {
	metaJSON, err := json.Marshal(metaData)
	if err != nil {
		log.Println(err)
	}
	uiHistory.Send(ui.Command{Action: "metaData", Args: []string{string(metaJSON)}})
}
