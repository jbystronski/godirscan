package navigator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/entry"
)

type Navigator struct {
	NumVisibleLines, DirSize, CurrentIndex int
	StartLine, EndLine                     int
	StartCell                              int
	RowWidth                               int
	SortMode                               uint8
	Entries                                []*entry.Entry
	backTrace                              []int
	CurrentPath, RootPath                  string
	ActiveRowIndex                         int
	ActiveRowText                          string
	IsActive                               bool
}

func NewNavigator() *Navigator {
	// StartLine := 0
	// EndLine := 0

	return &Navigator{
		// StartLine: &StartLine,
		// EndLine:   &EndLine,
		// NumVisibleLines: terminal.GetNumVisibleLines(),
	}
}

func (n *Navigator) AddBackTrace(index int) {
	n.backTrace = append(n.backTrace, index)
}

func (n *Navigator) GetBackTrace() int {
	if len(n.backTrace) == 0 {
		return 0
	}

	last := n.backTrace[len(n.backTrace)-1]
	n.backTrace = n.backTrace[:len(n.backTrace)-1]
	if last > n.GetEntriesLength() {
		return 0
	}
	return last
}

func (n *Navigator) HasEntries() bool {
	return n.GetEntriesLength() > 0
}

func (n *Navigator) GetEndLine() int {
	return n.EndLine
}

func (n *Navigator) GetCurrentIndex() int {
	return n.CurrentIndex
}

func (n *Navigator) SetCurrentIndex(i int) {
	n.CurrentIndex = i
}

func (n *Navigator) SetEndLine(i int) {
	n.EndLine = i
}

func (n *Navigator) GetStartLine() int {
	return n.StartLine
}

func (n *Navigator) SetStartLine(i int) {
	n.StartLine = i
}

func (n *Navigator) GetDirSize() *int {
	return &n.DirSize
}

func (n *Navigator) SetDirSize(num int) {
	n.DirSize = num
}

func (n *Navigator) SetCurrentPath(p string) {
	n.CurrentPath = p
}

func (n *Navigator) GetEntries() []*entry.Entry {
	return n.Entries
}

func (n *Navigator) SetEntries(newEntries []*entry.Entry) {
	n.Entries = newEntries
}

func (n *Navigator) SetRootPath(p string) {
	n.RootPath = p
}

func (n *Navigator) GetEntriesLength() int {
	return len(n.Entries)
}

func (n *Navigator) GetCurrentEntry() *entry.Entry {
	return n.Entries[n.CurrentIndex]
}

func (n *Navigator) GetEntry(index int) *entry.Entry {
	return n.Entries[index]
}

func (n *Navigator) Reset() {
	n.SetCurrentIndex(0)
	n.SetStartLine(0)
	n.SetEndLine(0)
}

func (n *Navigator) MoveDown() (ok bool) {
	if n.GetEntriesLength() > 0 && n.CurrentIndex < n.GetEntriesLength()-1 {
		n.CurrentIndex++
		if n.CurrentIndex >= n.StartLine+n.NumVisibleLines {
			n.StartLine++
		}
		ok = true

	}
	return
}

func (n *Navigator) MoveUp() (ok bool) {
	if n.HasEntries() && n.CurrentIndex > 0 {
		n.CurrentIndex--
		if n.CurrentIndex < n.StartLine {
			n.StartLine--
		}

		ok = true
	}
	return
}

func (n *Navigator) ClearEntries() {
	n.Entries = nil
}

func (n *Navigator) GetParentPath() string {
	parentDir, _ := filepath.Split(n.CurrentPath)

	if parentDir == string(os.PathSeparator) {
		return parentDir
	}

	return strings.TrimSuffix(parentDir, string(os.PathSeparator))
}

func (n *Navigator) IncrementIndex() {
	if n.CurrentIndex < n.GetEntriesLength()-1 {
		n.CurrentIndex++
	}
}
