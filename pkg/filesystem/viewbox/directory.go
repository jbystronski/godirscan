package viewbox

func (v FsViewBox) Directory(sep, line string) string {
	return v.BuildString(bold, sep, v.Theme().Main, line, reset)
}
