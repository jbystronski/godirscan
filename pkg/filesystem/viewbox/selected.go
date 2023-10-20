package viewbox

func (v FsViewBox) SelectedRow(sep, line string) string {
	return v.BuildString(v.Theme().BgSelect, v.Theme().Select, bold, sep, line, reset)
}
