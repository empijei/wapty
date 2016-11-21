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
	Cookies     string
	Time        string
	/*
		✓ #
		✓ Host
		✓ Method (GET,POST,etc)
		✓ URL (actually, only the path)
		✓ Params (bool)
		✓ Edited (bool)
		✓ StatusCode
		✓ Length
		✓ MIME type (parse content-type response header)
		✓ Extension
		Title (maybe not?)
		Comment (user-defined)
		✓ SSL (just check if Response.TLS != nil)
		✓ IP (dig it)
		✓ Cookies
		✓ Time
	*/
}

func newReqRespMetaData(Id uint) *ReqRespMetaData {
	return &ReqRespMetaData{Id: Id}
}

//DISCLAIMER use original req AFTER editing the new one
func (rr *ReqResp) parseRequest(req *http.Request) {
	this := rr.MetaData
	this.Host = req.Host
	this.Method = req.Method
	this.Path = req.URL.Path
	if len(req.Form) == 0 {
		_ = req.ParseForm()
	}
	this.Params = len(req.Form) > 0
	status.RLock()
	this.Edited = status.ReqResps[this.Id].RawEditedReq != nil
	status.RUnlock()
	tmp := strings.Split(this.Path, ".")
	if !strings.Contains(tmp[len(tmp)-1], "/") {
		this.Extension = tmp[len(tmp)-1]
	}
	ips, err := net.LookupHost(strings.Split(this.Host, ":")[0])
	if err == nil && len(ips) >= 1 {
		this.IP = ips[0]
	}
	this.Time = time.Now().String()
	sendMetaData(this)
}

//DISCLAIMER use original res AFTER editing the new one
func (rr *ReqResp) parseResponse(res *http.Response) {
	this := rr.MetaData
	if !this.Edited {
		status.RLock()
		this.Edited = status.ReqResps[this.Id].RawEditedRes != nil
		status.RUnlock()
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
	ui.Send(ui.Command{Channel: HISTORYCHANNEL, Action: "metaData", Args: []string{string(metaJSON)}})
}
