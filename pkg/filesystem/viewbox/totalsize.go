package viewbox

func (v FsViewBox) TotalSize(s int) string {
	return v.BuildString(v.Theme().BgHeader, v.Theme().Header, v.PrintSizeAsString(s), reset)
}
