package config

import t "github.com/jbystronski/godirscan/pkg/lib/termui"

type Color string

type Theme struct {
	Main, Accent, BgHighlight, Highlight, BgHeader, Header, BgSelect, Select, BgPrompt, Prompt string
}

type Schema struct {
	Main        Color `json:"fg_main"`
	Accent      Color `json:"fg_accent"`
	Highlight   Color `json:"fg_highlight"`
	Header      Color `json:"fg_header"`
	Select      Color `json:"fg_selected"`
	Prompt      Color `json:"fg_prompt"`
	BgHighlight Color `json:"bg_highlight"`
	BgHeader    Color `json:"bg_header"`
	BgSelect    Color `json:"bg_selected"`
	BgPrompt    Color `json:"bg_prompt"`
}

const (
	magenta  Color = "magenta"
	bmagenta Color = "bright_magenta"
	red      Color = "red"
	bred     Color = "bright_red"
	cyan     Color = "cyan"
	bcyan    Color = "bright_cyan"
	yellow   Color = "yellow"
	byellow  Color = "bright_yellow"
	black    Color = "black"
	bblack   Color = "bright_black"
	white    Color = "white"
	bwhite   Color = "bright_white"
	blue     Color = "blue"
	bblue    Color = "bright_blue"
	green    Color = "green"
	bgreen   Color = "bright_green"
)

var background = map[Color]string{
	black:    "\033[40m",
	red:      "\033[41m",
	green:    "\033[42m",
	yellow:   "\033[43m",
	blue:     "\033[44m",
	magenta:  "\033[45m",
	cyan:     "\033[46m",
	white:    "\033[47m",
	bblack:   "\033[100m",
	bred:     "\033[101m",
	bgreen:   "\033[102m",
	byellow:  "\033[103m",
	bblue:    "\033[104m",
	bmagenta: "\033[105m",
	bcyan:    "\033[106m",
	bwhite:   "\033[107m",
}

var foreground = map[Color]string{
	black:    "\033[30m",
	red:      "\033[31m",
	green:    "\033[32m",
	yellow:   "\033[33m",
	blue:     "\033[34m",
	magenta:  "\033[35m",
	cyan:     "\033[36m",
	white:    "\033[37m",
	bblack:   "\033[90m",
	bred:     "\033[91m",
	bgreen:   "\033[92m",
	byellow:  "\033[93m",
	bblue:    "\033[94m",
	bmagenta: "\033[95m",
	bcyan:    "\033[96m",
	bwhite:   "\033[97m",
}

/*

prebuilt schemas

*/

var PrebuiltSchemas = []Schema{
	{
		Main:        t.Magenta,
		Accent:      t.BrightCyan,
		Highlight:   t.Black,
		BgHighlight: t.BrightWhite,
		BgHeader:    t.Cyan,
		Header:      t.Black,
		Select:      t.BrightWhite,
		BgSelect:    t.Magenta,
		Prompt:      t.BrightWhite,
		BgPrompt:    t.Cyan,
	},
	{
		Main:        t.BrightBlue,
		Accent:      t.BrightYellow,
		BgHighlight: t.BrightWhite,
		Highlight:   t.Black,
		BgHeader:    t.BrightBlue,
		Header:      t.BrightWhite,
		Select:      t.Blue,
		BgSelect:    t.BrightYellow,
		Prompt:      t.BrightWhite,
		BgPrompt:    t.Yellow,
	},
	{
		Main:        t.BrightBlack,
		Accent:      t.BrightWhite,
		BgHighlight: t.BrightWhite,
		Highlight:   t.Black,
		Header:      t.Black,
		BgHeader:    t.BrightBlack,
		Select:      t.BrightWhite,
		BgSelect:    t.BrightBlack,
		Prompt:      t.BrightWhite,
		BgPrompt:    t.BrightBlack,
	},
	{
		Main:        t.Red,
		Accent:      t.BrightYellow,
		BgSelect:    t.BrightYellow,
		Select:      t.Black,
		BgHighlight: t.BrightWhite,
		Highlight:   t.Black,
		Header:      t.BrightYellow,
		BgHeader:    t.BrightRed,
		Prompt:      t.Red,
		BgPrompt:    t.BrightYellow,
	},
	{
		Main:        t.BrightCyan,
		Accent:      t.White,
		BgSelect:    t.BrightBlack,
		Select:      t.BrightWhite,
		BgHighlight: t.Cyan,
		Highlight:   t.White,
		Header:      t.BrightWhite,
		BgHeader:    t.Cyan,
		Prompt:      t.BrightWhite,
		BgPrompt:    t.Cyan,
	},
	{
		Main:        t.Yellow,
		Accent:      t.White,
		BgSelect:    t.Yellow,
		Select:      t.Black,
		BgHighlight: t.White,
		Highlight:   t.Black,
		Header:      t.BrightWhite,
		BgHeader:    t.BrightBlack,
		Prompt:      t.BrightWhite,
		BgPrompt:    t.Magenta,
	},
	{
		Main:        t.BrightGreen,
		Accent:      t.BrightWhite,
		Highlight:   t.Black,
		Header:      t.BrightWhite,
		Select:      t.BrightWhite,
		Prompt:      t.BrightWhite,
		BgHighlight: t.BrightWhite,
		BgHeader:    t.BrightMagenta,
		BgSelect:    t.BrightMagenta,
		BgPrompt:    t.BrightMagenta,
	},
}
