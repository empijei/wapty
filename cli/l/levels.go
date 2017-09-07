package l

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
	repr[Level_Debug],
	printColor(repr[Level_Info], color_FgCyan),
	printColor(repr[Level_Warn], color_FgYellow),
	printColor(repr[Level_Error], color_FgRed),
	printColor(repr[Level_Failure], color_FgMagenta),
}

func init() {
	for i := 0; i < len(reprC); i++ {
		reprC[i] = printColor(reprC[i], style_Bold)
	}
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
