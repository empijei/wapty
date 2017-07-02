package apis

// Action is a string representing the action to perform
type Action string

const (
	// Editor possible user actions

	// FORWARD tells the backend to forward the currently intercepted request/response
	FORWARD Action = "forward"
	// EDIT tells the backend to forward the provided payload instead of the original one
	EDIT = "edit"
	// DROP tells the backend to create a dummy response and send it to the client.
	// If a request is dropped it won't be forwarded to the server.
	DROP = "drop"
	// PROVIDERESP only has meaning if a request was intercepted. It allows to
	// provide a response to the current request without forwarding it to the server.
	PROVIDERESP = "provideResp"

	//History actions

	// DUMP dumps the entire status, this is for debug purposes only
	DUMP = "dump"
	// FILTER allows to search and filter through history
	FILTER = "filter"
	// FETCH return the ReqResp given an ID
	FETCH = "fetch"
	// METADATA returns the metadata for the given ID
	METADATA = "metaData"

	//Repeater actions

	// CREATE Creates a new repeater entry
	CREATE = "create"
	// GO Performs the request
	GO = "go"
	// GET Retrieves an history item
	GET = "get"

	//Settings actions

	// INTERCEPT Gets (no params) or sets (param "ON") the intercept status
	INTERCEPT = "intercept"
)

const (
	//Possible payloads types

	// REQUEST should be used to tell the UI if the payload is a request
	REQUEST = "request"
	// RESPONSE should be used to tell the UI if the payload is a response
	RESPONSE = "response"
)

// ArgName is the type of the set of Keys to use in the Args map of a command
type ArgName string

const (
	// ID is used to identify which item of the collection should be fetched
	ID ArgName = "id"
	// SUBID is used to identify which item of the collection should be fetched.
	// This is applied only if the resource identified by ID contains a collection.
	SUBID = "subId"
	// PAYLOADTYPE is used by editor to distinguish between requests and responses.
	PAYLOADTYPE = "payloadType"
	// ENDPOINT is used to refer to the host:port or schema://host:port
	ENDPOINT = "endpoint"
	// ERR is a value used to communicate an error occourred
	ERR = "error"
	// TLS is used as a bool to tell if TLS must be used for the specified operation
	TLS = "tls"
	// TRUE is used to deserialize a bool from a string.
	TRUE = "true"
	// FALSE is used to deserialize a bool from a string.
	FALSE = ""
	// ON is used as a key value for togglable settings
	ON = "on"
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
