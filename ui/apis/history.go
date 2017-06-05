package apis

//Enum for possible user actions
type HistoryAction int

const (
	DUMP HistoryAction = iota
	FILTER
	FETCH
	METADATA
)

var historyActions = [...]string{
	"dump",
	"filter",
	"fetch",
	"metaData",
}

func (a HistoryAction) String() string {
	return historyActions[a]
}

type ReqRespMetaData struct {
	Id          int
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

func NewReqRespMetaData(Id int) *ReqRespMetaData {
	return &ReqRespMetaData{Id: Id}
}

//Represents an item of the proxy history
//TODO create a test that fails if this is different from intercept.ReqResp
type ReqResp struct {
	//Unique Id in the history
	Id int
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
