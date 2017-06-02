package apis

type UiChannel int

const (
	EDITORCHANNEL UiChannel = iota
	HISTORYCHANNEL
	REPEATCHANNEL
	SETTINGSCHANNEL
)

var UiChannels = [...]string{
	"proxy/intercept/editor",
	"proxy/httpHistory",
	"repeat",
	"proxy/intercept/options",
}

func (uic UiChannel) String() string {
	return UiChannels[uic]
}
