package apis

//String used to recognize commands directed to this module
const HISTORYCHANNEL = "proxy/httpHistory"

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

func NewReqRespMetaData(Id uint) *ReqRespMetaData {
	return &ReqRespMetaData{Id: Id}
}
