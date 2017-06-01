package apis

//String used to recognize commands directed to this module
const HISTORYCHANNEL = "proxy/httpHistory"

//Enum for possible user actions
type HistoryAction int

const (
	DUMP HistoryAction = iota
	FILTER
	FETCH
)

var historyActions = [...]string{
	"dump",
	"filter",
	"fetch",
}

func (a HistoryAction) String() string {
	return historyActions[a]
}
