package apis

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
