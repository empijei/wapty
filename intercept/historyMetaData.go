package intercept

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/empijei/wapty/ui/apis"
)

//DISCLAIMER use original req AFTER editing the new one
//And use it from a thread that has a readlock on the status
func (rr *ReqResp) parseRequest(req *http.Request) {
	this := rr.MetaData
	this.Host = req.Host
	this.Method = req.Method
	this.Path = req.URL.Path
	//TODO implement this in a way that does not consume the body
	//if len(req.Form) == 0 {
	//	_ = req.ParseForm()
	//}
	//this.Params = len(req.Form) > 0
	//this supposes to alread have a RLock on the status.
	this.Edited = status.ReqResps[this.Id].RawEditedReq != nil
	tmp := strings.Split(this.Path, ".")
	if !strings.Contains(tmp[len(tmp)-1], "/") {
		this.Extension = tmp[len(tmp)-1]
	}
	ipport := strings.Split(this.Host, ":")
	ips, err := net.LookupHost(ipport[0])
	if err == nil && len(ips) >= 1 {
		this.IP = ips[0]
		if len(ipport) >= 2 {
			this.Port = ipport[1]
		} else {
			switch req.URL.Scheme {
			case "https":
				this.Port = "443"
			case "http":
				this.Port = "80"
			default:
				log.Println("Port not specified: " + this.Host)
			}
		}
	} else {
		log.Println("Unable to resolve Host: " + this.Host)
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
	//FIXME
	this.Length = res.ContentLength
	this.ContentType = res.Header.Get("Content-Type")
	this.TLS = res.TLS != nil
	tmp := res.Cookies()
	for _, cookie := range tmp {
		this.Cookies += cookie.String() + "; "
	}
	sendMetaData(this)
}

func sendMetaData(metaData *apis.ReqRespMetaData) {
	metaJSON, err := json.Marshal(metaData)
	if err != nil {
		log.Println(err)
	}
	uiHistory.Send(apis.Command{Action: apis.METADATA, Args: map[apis.Param]string{apis.METADATA: string(metaJSON)}})
}
