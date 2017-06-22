package apis

const (
	ID          = "id"
	SUBID       = "subId"
	PAYLOADTYPE = "payloadType"
	ENDPOINT    = "endpoint"
	ERR         = "error"
	TLS         = "tls"
	HOST        = "host"
	TRUE        = "true"
	FALSE       = ""
	ON          = "on"
)

func Err(message string) *Command {
	return &Command{
		Action: ERR,
		Args:   map[string]string{ERR: message},
	}
}
