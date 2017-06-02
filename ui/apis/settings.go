package apis

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
