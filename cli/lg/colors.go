package lg

import "fmt"

// Thanks to "github.com/fatih/color" for these constants

type cliAttribute int

const escape = "\x1b"

// Base attributes
const (
	style_Reset cliAttribute = iota
	style_Bold
	style_Faint
	style_Italic
	style_Underline
	style_BlinkSlow
	style_BlinkRapid
	style_ReverseVideo
	style_Concealed
	style_CrossedOut
)

// Foreground text colors
const (
	color_FgBlack cliAttribute = iota + 30
	color_FgRed
	color_FgGreen
	color_FgYellow
	color_FgBlue
	color_FgMagenta
	color_FgCyan
	color_FgWhite
)

// Foreground Hi-Intensity text colors
const (
	color_FgHiBlack cliAttribute = iota + 90
	color_FgHiRed
	color_FgHiGreen
	color_FgHiYellow
	color_FgHiBlue
	color_FgHiMagenta
	color_FgHiCyan
	color_FgHiWhite
)

// Background text colors
const (
	color_BgBlack cliAttribute = iota + 40
	color_BgRed
	color_BgGreen
	color_BgYellow
	color_BgBlue
	color_BgMagenta
	color_BgCyan
	color_BgWhite
)

// Background Hi-Intensity text colors
const (
	color_BgHiBlack cliAttribute = iota + 100
	color_BgHiRed
	color_BgHiGreen
	color_BgHiYellow
	color_BgHiBlue
	color_BgHiMagenta
	color_BgHiCyan
	color_BgHiWhite
)

func printColor(s string, color cliAttribute) string {
	return fmt.Sprintf("%s[%dm%s%s[%dm", escape, color, s, escape, style_Reset)
}

func printColorStyle(s string, color cliAttribute, style cliAttribute) string {
	return fmt.Sprintf("%s[%d;%dm%s%s[%dm", escape, style, color, s, escape, style_Reset)
}
