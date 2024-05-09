package filesystem

import (
	"strings"

	"github.com/jbystronski/godirscan/pkg/app/data"
)

func (c *FsController) filter(search string) bool {
	result := []*data.FsEntry{}

	for _, en := range c.data.All() {
		if strings.Contains(strings.ToLower(en.Name()), strings.ToLower(search)) {
			result = append(result, en)
		}
	}

	if len(result) == 0 {
		return false
	}

	c.data.SetData(result)
	return true
}
