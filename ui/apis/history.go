package apis

// ReqRespMetaData is a wrapper type to hold all metadata on a status ReqResp
type ReqRespMetaData struct {
	ID          int
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

// ReqResp Represents an item of the proxy history
//TODO create a test that fails if this is different from intercept.ReqResp
type ReqResp struct {
	//Unique ID in the history
	ID int
	//Meta Data about both Req and Resp
	MetaData *ReqRespMetaData
	//Original Request
	RawReq []byte
	//Original Response
	RawRes []byte
	//Edited Request
	RawEditedReq []byte
	//Edited Response
	RawEditedRes []byte
}
