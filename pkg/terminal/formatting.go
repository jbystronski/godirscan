package terminal

import "strings"

const (
	Space = " "

	Line = "\u2594"

	Bslash     = "\u2572"
	Fslash     = "\u2571"
	DbVert     = "\u2551"
	Texture    = "\u2591"
	ResetFmt   = "\033[0m"
	BoldFmt    = "\033[1m"
	DbHoriz    = "\u2550"
	CursorTop  = "\033[H"
	Hseparator = "\u2500"
	CornerLine = "\u2514"
	TeeLine    = "\u251c"
)

var Segment = strings.Join([]string{Fslash, Texture, Fslash}, "")

var BackgroundsMap = map[string]string{
	"black":          "\033[40m",
	"red":            "\033[41m",
	"green":          "\033[42m",
	"yellow":         "\033[43m",
	"blue":           "\033[44m",
	"magenta":        "\033[45m",
	"cyan":           "\033[46m",
	"white":          "\033[47m",
	"bright_black":   "\033[100m",
	"bright_red":     "\033[101m",
	"bright_green":   "\033[102m",
	"bright_yellow":  "\033[103m",
	"bright_blue":    "\033[104m",
	"bright_magenta": "\033[105m",
	"bright_cyan":    "\033[106m",
	"bright_white":   "\033[107m",
}

var ColorsMap = map[string]string{
	"black":          "\033[30m",
	"red":            "\033[31m",
	"green":          "\033[32m",
	"yellow":         "\033[33m",
	"blue":           "\033[34m",
	"magenta":        "\033[35m",
	"cyan":           "\033[36m",
	"white":          "\033[37m",
	"bright_black":   "\033[90m",
	"bright_red":     "\033[91m",
	"bright_green":   "\033[92m",
	"bright_yellow":  "\033[93m",
	"bright_blue":    "\033[94m",
	"bright_magenta": "\033[95m",
	"bright_cyan":    "\033[96m",
	"bright_white":   "\033[97m",
}

// const (
// 	EmptyIndent    = "   "
// 	VerticalSep    = "\u2502"
//
//
// 	HorizontalLine = "\u2500"
//
//

// 	FmtDir         = Bold + Blue
// 	FmtFile        = Yellow
// )

// boldFmt    = "\033[1m"
// italic     = "\033[3m"
// resetFmt   = "\033[0m"
// underscore = "\033[4m"
// black      = "\033[30m"

// tlCorner      = "\033(0\x6C\033(B"
// trCorner      = "\033(0\x71\033(B"
// blCorner      = "\033(0\x6D\033(B"
// brCorner      = "\033(0\x6A\033(B"
// lrDiag        = "\033(0\x2F\033(B"
// hTop          = "\033(0\x48\033(B"

// emptyIndent = "   "
// vSeparator  = "\u2502"
// teeLine     = "\u251c"
// cornerLine  = "\u2514"
// hSeparator  = "\u2500"
// dbHoriz     = "\u2550"
// line        = "\u2594"
// dbVert      = "\u2551"
// segment     = "\u2571" + "\u2591" + "\u2571"
// block       = "\u2586"
// space       = " "
// fSlash      = "\u2571"
// bSlash      = "\u2572"
//
// texture     = "\u2591"
