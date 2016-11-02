package intercept

import "github.com/empijei/Wapty/ui"

//String used to recognize commands directed to this module
const HISTORYCHANNEL = "proxy/httpHistory"

var uiHistory *ui.Subscription

//Enum for possible user actions
type HistoryAction int

const (
	FETCH HistoryAction = iota
	FILTER
)

var historyActions = [...]string{
	"fetch",
	"filter",
}

func (a HistoryAction) String() string {
	return historyActions[a]
}

var invertHistoryActions map[string]HistoryAction

func parseAction(s string) HistoryAction {
	return invertHistoryActions[s]
}
