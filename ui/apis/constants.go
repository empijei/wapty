package apis

// Action is a string representing the action to perform
type Action string

const (
	// Editor possible user actions

	// EDT_FORWARD tells the backend to forward the currently intercepted request/response
	EDT_FORWARD Action = "forward"
	// EDT_EDIT tells the backend to forward the provided payload instead of the original one
	EDT_EDIT = "edit"
	// EDT_DROP tells the backend to create a dummy response and send it to the client.
	// If a request is dropped it won't be forwarded to the server.
	EDT_DROP = "drop"
	// EDT_PROVIDERESP only has meaning if a request was intercepted. It allows to
	// provide a response to the current request without forwarding it to the server.
	EDT_PROVIDERESP = "provideResp"

	//History actions

	// HST_DUMP dumps the entire status, this is for debug purposes only
	HST_DUMP = "dump"
	// HST_FILTER allows to search and filter through history
	HST_FILTER = "filter"
	// HST_FETCH return the ReqResp given an ID
	HST_FETCH = "fetch"
	// HST_METADATA returns the metadata for the given ID
	HST_METADATA = "metaData"

	//Repeater actions

	// RPT_CREATE Creates a new repeater entry
	RPT_CREATE = "create"
	// RPT_GO Performs the request
	RPT_GO = "go"
	// RPT_GET Retrieves an history item
	RPT_GET = "get"

	//Settings actions

	// STN_INTERCEPT Gets (no params) or sets (param "ON") the intercept status
	STN_INTERCEPT = "intercept"
)

const (
	//Possible payloads types

	// PLD_REQUEST should be used to tell the UI if the payload is a request
	PLD_REQUEST = "request"
	// PLD_RESPONSE should be used to tell the UI if the payload is a response
	PLD_RESPONSE = "response"
)

// ArgName is the type of the set of Keys to use in the Args map of a command
type ArgName string

const (
	// ARG_ID is used to identify which item of the collection should be fetched
	ARG_ID ArgName = "id"
	// ARG_SUBID is used to identify which item of the collection should be fetched.
	// This is applied only if the resource identified by ID contains a collection.
	ARG_SUBID = "subId"
	// ARG_PAYLOADTYPE is used by editor to distinguish between requests and responses.
	ARG_PAYLOADTYPE = "payloadType"
	// ARG_ENDPOINT is used to refer to the host:port or schema://host:port
	ARG_ENDPOINT = "endpoint"
	// ARG_ERR is a value used to communicate an error occourred
	ARG_ERR = "error"
	// ARG_TLS is used as a bool to tell if ARG_TLS must be used for the specified operation
	ARG_TLS = "tls"
	// ARG_TRUE is used to deserialize a bool from a string.
	ARG_TRUE = "true"
	// ARG_FALSE is used to deserialize a bool from a string.
	ARG_FALSE = ""
	// ARG_ON is used as a key value for togglable settings
	ARG_ON = "on"
)

// UIChannel is a string used to multiplex on the websocket and route commands
// to the proper packages
type UIChannel string

const (
	// EDITORCHANNEL channel used by intercept package, editor actions
	EDITORCHANNEL UIChannel = "proxy/intercept/editor"
	// HISTORYCHANNEL channel used by intercept package, history actions
	HISTORYCHANNEL = "proxy/httpHistory"
	// REPEATCHANNEL channel used by repeat package
	REPEATCHANNEL = "repeat"
	// INTERCEPTSETTINGSCHANNEL channel used by intercept package, history actions
	INTERCEPTSETTINGSCHANNEL = "proxy/intercept/options"
)
