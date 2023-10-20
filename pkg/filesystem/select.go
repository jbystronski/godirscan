package filesystem

import (
	"path/filepath"

	"github.com/jbystronski/godirscan/pkg/common"
)

type Selected struct {
	common.MapAccessor[string, struct{}]
}

func NewSelected() *Selected {
	return &Selected{
		MapAccessor: &common.GenericMap[string, struct{}]{},
	}
}

func (s *Selected) Toggle(path string) {
	if s.Exists(path) {
		s.Unset(path)
		return
	}
	s.Set(path, struct{}{})
}

func (s Selected) BasePath() string {
	if s.Len() == 0 {
		return ""
	}

	for k := range s.Self() {

		dir, _ := filepath.Split(k)
		return dir

	}

	return ""
}
