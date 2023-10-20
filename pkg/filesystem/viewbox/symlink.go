package viewbox

func (v FsViewBox) Symlink(sep, line string) string {
	return v.BuildString(bold, sep, v.Theme().Main, line, reset)
}
