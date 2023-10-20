package viewbox

import "fmt"

func (v FsViewBox) PrintEmpty(isPanelActive bool) {
	for i := v.OutputFirstLine(); i <= v.OutputLastLine(); i++ {
		v.ClearLine(i, v.ContentLineStart(), v.ContentWidth())
	}

	v.GoToCell(v.OutputFirstLine(), v.ContentLineStart())

	var output string
	if isPanelActive {
		output = v.BuildString(v.Theme().BgHighlight, v.Theme().Highlight, bold, "Folder is empty", reset)
	} else {
		output = v.BuildString("Folder is empty")
	}

	fmt.Print(output)
}
