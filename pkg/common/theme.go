package common

type Theme struct {
	Main        string `json:"main_color"`
	Accent      string `json:"accent_color"`
	BgHighlight string `json:"highlight_background"`
	Highlight   string `json:"highligh_color"`
	BgHeader    string `json:"header_background"`
	Header      string `json:"header_color"`
	BgSelect    string `json:"selected_background"`
	Select      string `json:"selected_color"`
	BgPrompt    string `json:"prompt_background"`
	Prompt      string `json:"prompt_color"`
}

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
