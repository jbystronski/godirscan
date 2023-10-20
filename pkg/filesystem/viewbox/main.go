package viewbox

import (
	"fmt"
	"strings"

	"github.com/jbystronski/godirscan/pkg/common"
	"github.com/jbystronski/godirscan/pkg/converter"
	"github.com/jbystronski/godirscan/pkg/viewbox"
)

const (
	horizontalBorder        = "\u2501"
	verticalBorder          = "\u2503"
	topLeftBorderCorner     = "\u250F"
	topRightBorderCorner    = "\u2513"
	bottomLeftBorderCorner  = "\u2517"
	bottomRightBorderCorner = "\u251B"
	horizontalSeparator     = "\u2500"
	corner                  = "\u2514"
	tee                     = "\u251c"
	reset                   = "\033[0m"
	bold                    = "\033[1m"
)

type FsViewBox struct {
	viewbox.ViewBox
}

func (v FsViewBox) GoToPromptCell() common.Coords {
	return common.Coords{
		Y: v.TotalLines() - 2,
		X: v.ContentLineStart(),
	}
}

func (v FsViewBox) GoToTotalSizeCell() common.Coords {
	return v.GoToPromptCell()
}

func (v FsViewBox) Line(index int) int {
	return (index % v.Lines()) + v.FirstLine()
}

func FormatSize(bytes int) string {
	if bytes < int(converter.KbInBytes) {
		st := fmt.Sprintf("%d %s", bytes, converter.StorageUnits[0]+" ")
		if len(st) < sizeMaxLen {
			st = strings.Repeat(" ", sizeMaxLen-len(st)) + st
		}
		return " " + st
	}

	floatSize, unit := converter.BytesToFloat(bytes)

	st := fmt.Sprintf("%.2f %s", floatSize, unit+" ")

	if len(st) < sizeMaxLen {
		st = strings.Repeat(" ", sizeMaxLen-len(st)) + st
	}
	return " " + st
}

const sizeMaxLen = 11
