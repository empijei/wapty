package apis

// Action is a string representing the action to perform
type Action string

const (
	//Editor possible user actions
	FORWARD     Action = "forward"
	EDIT               = "edit"
	DROP               = "drop"
	PROVIDERESP        = "provideResp"

	//Possible payloads types
	REQUEST  = "request"
	RESPONSE = "response"

	//History actions
	DUMP     = "dump"
	FILTER   = "filter"
	FETCH    = "fetch"
	METADATA = "metaData"

	//Repeater actions
	//Creates a new repeater entry
	CREATE = "create"
	//Performs the request
	GO = "go"
	//Retrieves an history item
	GET = "get"

	//Settings actions
	//Gets (no params) or sets (param "ON") the intercept status
	INTERCEPT = "intercept"
)

// ArgName is the type of the set of Keys to use in the Args map of a command
type ArgName string

const (
	ID          ArgName = "id"
	SUBID               = "subId"
	PAYLOADTYPE         = "payloadType"
	ENDPOINT            = "endpoint"
	ERR                 = "error"
	TLS                 = "tls"
	TRUE                = "true"
	FALSE               = ""
	ON                  = "on"
)

// UIChannel is a string used to multiplex on the websocket and route commands
// to the proper packages
type UIChannel string

const (
	EDITORCHANNEL   UIChannel = "proxy/intercept/editor"
	HISTORYCHANNEL            = "proxy/httpHistory"
	REPEATCHANNEL             = "repeat"
	SETTINGSCHANNEL           = "proxy/intercept/options"
)

type Command struct {
	Channel UIChannel
	Action  Action
	Args    map[ArgName]string
	Payload []byte
}

func Err(message string) *Command {
	return &Command{
		Action: ERR,
		Args:   map[ArgName]string{ERR: message},
	}
}
