package viewbox

func (v FsViewBox) ActiveRow(sep, line string) string {
	return v.BuildString(v.Theme().BgHighlight, v.Theme().Highlight, bold, sep, line, reset)
}
