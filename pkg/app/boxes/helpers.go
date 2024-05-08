package boxes

import (
	g "github.com/jbystronski/godirscan/pkg/global"
	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

var (
	updateDimensions = g.UpdateDimensions
	alignRight       = g.AlignRight
	alignLeft        = g.AlignLeft
	alignCenter      = g.AlignCenter
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
	fmtSize          = g.FormatSize
)

const Reset = termui.Reset
