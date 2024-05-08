package app

import (
	"github.com/jbystronski/godirscan/pkg/app/config"
	g "github.com/jbystronski/godirscan/pkg/global"
)

var (
	updateDimensions = g.UpdateDimensions
	hideCursor       = g.HideCurson
	showCursor       = g.ShowCursor
	strlen           = g.StrLen
	cls              = g.ClearScreen
	cell             = g.Cell
	clear            = g.Clear
	rows             = g.Rows
	cols             = g.Cols
	trimEnd          = g.TrimEnd
	buildString      = g.BuildString
	fmtBold          = g.FmtBold
	themeMain        = g.ThemeMain
	ThemeAccent      = g.ThemeAccent
	themeBgHighlight = g.ThemeBgHighlight
	themeHighlight   = g.ThemeHighlight
	themeBgHeader    = g.ThemeBgHeader
	themeHeader      = g.ThemeHeader
	themeBgSelect    = g.ThemeBgSelect
	themeSelect      = g.ThemeSelect
	fmtPrompt        = g.FmtPrompt
)

const (
	MIN_HEIGHT = config.MIN_HEIGHT
	MIN_WIDTH  = config.MIN_WIDTH
)
