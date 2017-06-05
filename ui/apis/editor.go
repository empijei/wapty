package apis

//Enum for possible user actions
type EditorAction int

const (
	FORWARD EditorAction = iota
	EDIT
	DROP
	PROVIDERESP
)

var editorActions = [...]string{
	"forward",
	"edit",
	"drop",
	"provideResp",
}

func (a EditorAction) String() string {
	return editorActions[a]
}

//Enum for possible payloads types
type PayloadType int

const (
	REQUEST PayloadType = iota
	RESPONSE
)

var payloads = [...]string{
	"request",
	"response",
}

func (p PayloadType) String() string {
	return payloads[p]
}
