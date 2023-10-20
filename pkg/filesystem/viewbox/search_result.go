package viewbox

import (
	"github.com/jbystronski/godirscan/pkg/filesystem"
)

func (v FsViewBox) SearchResult(sep string, en filesystem.FsEntry) string {
	return v.BuildString(bold, sep, v.Theme().Accent, en.FullPath(), reset)
}
