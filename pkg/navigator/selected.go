package navigator

import "github.com/jbystronski/godirscan/pkg/entry"

type Selected struct {
	SelectedEntries     map[*entry.Entry]struct{}
	selectedEntriesPath string
}

func NewSelected() *Selected {
	return &Selected{
		SelectedEntries: make(map[*entry.Entry]struct{}),
	}
}

func (s *Selected) Clear() {
	s.SelectedEntries = make(map[*entry.Entry]struct{})
}

func (s *Selected) Select(e *entry.Entry) {
	if _, ok := s.SelectedEntries[e]; ok {
		delete(s.SelectedEntries, e)
	} else {
		s.SelectedEntries[e] = struct{}{}
	}
}

func (s *Selected) SelectAll(entries []*entry.Entry) {
	if len(s.SelectedEntries) > 0 {
		s.Clear()
	} else {
		for _, entry := range entries {
			s.SelectedEntries[entry] = struct{}{}
		}
	}
}

func (s *Selected) DumpPrevious(path string) {
	if path != s.selectedEntriesPath {
		s.Clear()
		s.selectedEntriesPath = path
	}
}

func (s *Selected) IsEmpty() bool {
	return len(s.SelectedEntries) == 0
}

func (s *Selected) GetAll() map[*entry.Entry]struct{} {
	return s.SelectedEntries
}
