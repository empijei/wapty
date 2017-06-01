package apis

//String used to recognize commands directed to this module
const SETTINGSCHANNEL = "proxy/intercept/options"

//Enum for possible user actions
type SettingsAction int

const (
	INTERCEPT SettingsAction = iota
)

var settingsActions = [...]string{
	"intercept",
}

func (a SettingsAction) String() string {
	return settingsActions[a]
}
