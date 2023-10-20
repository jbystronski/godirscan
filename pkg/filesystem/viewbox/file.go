package viewbox

func (v FsViewBox) File(sep, line string) string {
	return v.BuildString(bold, sep, v.Theme().Accent, line, reset)
}
