package filesystem

import (
	"fmt"
	"strings"

	"github.com/jbystronski/godirscan/pkg/common"
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

type ViewBox struct {
	common.ViewBox
	common.Trimmer
}

type Entity FsFiletype

type UpdateLineOptions struct {
	en    FsFiletype
	index int
	isActive,
	isControllerActive,
	isSelected,
	isLast bool
}

type RenderOpts struct {
	start,
	end,
	storeSize,
	index int
	entries            []*FsFiletype
	checkSelected      func(string) bool
	selectKey          func(FsFiletype) string
	storeName          string
	isControllerActive bool
}

func (v *ViewBox) printActiveRow(separator string, en fmt.Stringer) string {
	return fmt.Sprintf("%s%s%s%s%s\n", v.Theme().BgHighlight, v.Theme().Highlight+bold, separator, en, reset)
}

func (v *ViewBox) printSelectedRow(sep string, en fmt.Stringer) string {
	return fmt.Sprintf("%s%s%s%s%s\n", v.Theme().BgSelect, v.Theme().Select+bold, sep, en, reset)
}

func (v *ViewBox) printDirectory(sep string, en fmt.Stringer) string {
	return fmt.Sprintf("%v%s%v%s%v", bold, sep, v.Theme().Main, en, reset)
}

func (v *ViewBox) printFile(sep string, en fmt.Stringer) string {
	return fmt.Sprintf("%v%s%v%s%v", bold, sep, v.Theme().Accent, en, reset)
}

func (v *ViewBox) printSymlink(sep string, en fmt.Stringer) string {
	return fmt.Sprintf("%v%s%v%s%v", bold, sep, v.Theme().Main, en, reset)
}

func (v *ViewBox) printHeader(h fmt.Stringer) {
	fmt.Printf("%s%s%s", v.Theme().BgHeader+v.Theme().Header, h, reset)
}

func (v *ViewBox) printEmptyFolder(isPanelActive bool) {
	for i := v.OutputFirstLine(); i <= v.OutputLastLine(); i++ {
		v.ClearLine(i, v.ContentLineStart(), v.ContentWidth())
	}

	v.GoToCell(v.OutputFirstLine(), v.ContentLineStart())

	if isPanelActive {
		fmt.Print(v.Theme().BgHighlight, v.Theme().Highlight+bold, "Folder is empty", reset)
	} else {
		fmt.Print("Folder is empty")
	}
}

func (v *ViewBox) GoToPromptCell() common.Coords {
	return common.Coords{
		Y: v.TotalLines() - 2,
		X: v.ContentLineStart(),
	}
}

func (v *ViewBox) GoToTotalSizeCell() common.Coords {
	return v.GoToPromptCell()
}

func (v *ViewBox) Line(index int) int {
	return (index % v.Lines()) + v.FirstLine()
}

func (v *ViewBox) Separator(isLast bool) string {
	if isLast {
		return "\u2514\u2500"
	}
	return "\u251c\u2500"
}

func (v *ViewBox) RenderSingle(opts UpdateLineOptions) {
	v.ClearLine(v.Line(opts.index), v.ContentLineStart(), v.ContentWidth())

	var output string
	sep := v.Separator(opts.isLast)

	t := v.TrimEnd(fmt.Sprint(opts.en), v.ContentWidth()-16, 2)

	switch true {

	case opts.isActive && opts.isControllerActive:
		output = v.printActiveRow(sep, t)

	case opts.isSelected:
		output = v.printSelectedRow(sep, t)

	default:
		switch opts.en.(type) {
		case *FsDirectory:
			output = v.printDirectory(sep, t)

		case *FsSymlink:
			output = v.printSymlink(sep, t)

		case *FsFile:
			output = v.printFile(sep, t)

		default:
			output = v.printFile(sep, t)
		}
	}

	fmt.Print(opts.en.printSize(), output)
}

func (v *ViewBox) RenderAll(opts RenderOpts) {
	v.PrintBox()
	v.ClearLine(1, v.ContentLineStart()-1, v.Width())

	header := v.TrimEnd(opts.storeName, v.Width()-1, 2)

	v.printHeader(header)

	if len(opts.entries) == 0 {
		v.printEmptyFolder(opts.isControllerActive)
		return
	}

	v.GoToCell(v.Line(opts.index), v.ContentLineStart())

	line := v.FirstLine()

	for i := opts.start; i <= opts.end; i++ {
		file := opts.entries[i]

		if file == nil {
			continue
		}

		isSelected := opts.checkSelected(opts.selectKey(*file))

		opts := UpdateLineOptions{
			en:                 *file,
			index:              i,
			isActive:           i == opts.index,
			isControllerActive: opts.isControllerActive,
			isSelected:         isSelected,
			isLast:             i == len(opts.entries)-1,
		}

		v.RenderSingle(opts)

		line++
		v.GoToCell(line, v.ContentLineStart())
	}

	for line <= v.OutputLastLine() {
		v.ClearLine(line, v.ContentLineStart(), v.ContentWidth())
		line++
	}

	v.printTotalSize(opts.storeSize)
}

func (v *ViewBox) printTotalSize(s int) {
	v.GoToCell(v.TotalLines()-2, v.ContentLineStart())
	fmt.Printf("%v%v%s%v", v.Theme().BgHeader, v.Theme().Header, printSizeAsString(s), reset)

	v.GoToCell(v.TotalLines(), 1)
}

func FormatSize(bytes int) string {
	if bytes < int(common.KbInBytes) {
		st := fmt.Sprintf("%d %s", bytes, common.StorageUnits[0]+" ")
		if len(st) < sizeMaxLen {
			st = strings.Repeat(" ", sizeMaxLen-len(st)) + st
		}
		return " " + st
	}

	floatSize, unit := common.BytesToFloat(bytes)

	st := fmt.Sprintf("%.2f %s", floatSize, unit+" ")

	if len(st) < sizeMaxLen {
		st = strings.Repeat(" ", sizeMaxLen-len(st)) + st
	}
	return " " + st
}

func printSizeAsString(size int) string {
	return fmt.Sprintf("%v", FormatSize(size))
}

const sizeMaxLen = 11
