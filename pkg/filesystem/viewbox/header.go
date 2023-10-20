package viewbox

func (v FsViewBox) Header(h string) string {
	return v.BuildString(v.Theme().BgHeader, v.Theme().Header, h, reset)
}
