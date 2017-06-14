package apis

type UIChannel int

const (
	EDITORCHANNEL UIChannel = iota
	HISTORYCHANNEL
	REPEATCHANNEL
	SETTINGSCHANNEL
)

var UIChannels = [...]string{
	"proxy/intercept/editor",
	"proxy/httpHistory",
	"repeat",
	"proxy/intercept/options",
}

func (uic UIChannel) String() string {
	return UIChannels[uic]
}
