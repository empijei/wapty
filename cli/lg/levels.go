package lg

type LogLevel uint

const (
	Level_Debug LogLevel = iota
	Level_Info
	Level_Warn
	Level_Error
	Level_Failure
)

func (ll LogLevel) String() string {
	if ll > Level_Failure {
		return "Logging level out of bound"
	}
	var repr = []string{
		"Debug",
		"Info",
		"Warn",
		"Error",
		"Failure",
	}
	return repr[ll]
}

var repr = []string{
	"D",
	"I",
	"W",
	"E",
	"F",
}
var reprC = []string{
	printColor(repr[Level_Debug], style_Bold),
	printColorStyle(repr[Level_Info], color_FgCyan, style_Bold),
	printColorStyle(repr[Level_Warn], color_FgYellow, style_Bold),
	printColorStyle(repr[Level_Error], color_FgRed, style_Bold),
	printColorStyle(repr[Level_Failure], color_FgMagenta, style_Bold),
}

func (ll LogLevel) ShortString(color bool) string {
	if ll > Level_Failure {
		return "Logging level out of bound"
	}
	if color {
		return reprC[ll]
	}
	return repr[ll]
}
