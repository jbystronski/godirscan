package filesystem

import (
	"fmt"
	"os"
	"sync"
)

type Selected struct {
	entries map[string]struct{}
	path    string
}

func NewSelected() *Selected {
	return &Selected{
		entries: make(map[string]struct{}),
	}
}

func (s *Selected) Clear() {
	s.entries = make(map[string]struct{})
}

func (s *Selected) Select(path string) {
	if _, ok := s.entries[path]; ok {
		delete(s.entries, path)
	} else {
		s.entries[path] = struct{}{}
	}
}

func (s *Selected) Path() string {
	return s.path
}

func (s *Selected) SelectAll(entries Entries) {
	if len(s.entries) > 0 {
		s.Clear()
	} else {
		for _, entry := range entries.All() {
			s.entries[(*entry).FullPath()] = struct{}{}
		}
	}
}

func (s *Selected) DumpPrevious(path string) {
	if path != s.path {
		s.Clear()
		s.path = path
	}
}

func (s *Selected) IsEmpty() bool {
	return len(s.entries) == 0
}

func (s *Selected) All() map[string]struct{} {
	return s.entries
}

func (s *Selected) IsSelected(key string) (ok bool) {
	if _, ok = s.All()[key]; ok {
		return ok
	}
	return
}

func (s *Selected) Delete(path string, messageChan chan<- string) (ok bool, err error) {
	s.DumpPrevious(path)

	var wg sync.WaitGroup

	for key := range s.entries {

		wg.Add(1)
		go func(key string) {
			defer func() {
				wg.Done()
				messageChan <- fmt.Sprint("Deleted ", key, " ")
			}()

			error := os.RemoveAll(key)
			if error != nil {
				err = error
				return
			}
		}(key)

	}
	wg.Wait()

	ok = true

	return
}
