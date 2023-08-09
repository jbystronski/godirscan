package navigator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/entry"
)

type Navigator struct {
	NumVisibleLines, dirSize, currentIndex int
	startLine, endLine                     *int

	SortMode              uint8
	entries               []*entry.Entry
	backTrace             []int
	currentPath, rootPath string
}

func NewNavigator() *Navigator {
	startLine := 0
	endLine := 0

	return &Navigator{
		startLine: &startLine,
		endLine:   &endLine,
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

func (n *Navigator) GetEndline() int {
	return *n.endLine
}

func (n *Navigator) GetCurrentIndex() int {
	return n.currentIndex
}

func (n *Navigator) SetCurrentIndex(i int) {
	n.currentIndex = i
}

func (n *Navigator) SetEndLine(i int) {
	*n.endLine = i
}

func (n *Navigator) GetStartLine() int {
	return *n.startLine
}

func (n *Navigator) SetStartLine(i int) {
	*n.startLine = i
}

func (n *Navigator) GetDirSize() *int {
	return &n.dirSize
}

func (n *Navigator) SetDirSize(num int) {
	n.dirSize = num
}

func (n *Navigator) GetCurrentPath() string {
	return n.currentPath
}

func (n *Navigator) GetRootPath() string {
	return n.rootPath
}

func (n *Navigator) SetCurrentPath(p string) {
	n.currentPath = p
}

func (n *Navigator) GetEntries() []*entry.Entry {
	return n.entries
}

func (n *Navigator) SetEntries(newEntries []*entry.Entry) {
	n.entries = newEntries
}

func (n *Navigator) SetRootPath(p string) {
	n.rootPath = p
}

func (n *Navigator) GetEntriesLength() int {
	return len(n.entries)
}

func (n *Navigator) GetCurrentEntry() *entry.Entry {
	return n.entries[n.currentIndex]
}

func (n *Navigator) GetEntry(index int) *entry.Entry {
	return n.entries[index]
}

func (n *Navigator) Reset() {
	n.SetCurrentIndex(0)
	n.SetStartLine(0)
	n.SetEndLine(0)
}

func (n *Navigator) MoveDown() (ok bool) {
	if n.GetEntriesLength() > 0 && n.currentIndex < n.GetEntriesLength()-1 {
		n.currentIndex++
		if n.currentIndex >= *n.startLine+n.NumVisibleLines {
			*n.startLine++
		}
		ok = true

	}
	return
}

func (n *Navigator) MoveUp() (ok bool) {
	if n.HasEntries() && n.currentIndex > 0 {
		n.currentIndex--
		if n.currentIndex < *n.startLine {
			*n.startLine--
		}

		ok = true
	}
	return
}

func (n *Navigator) ClearEntries() {
	n.entries = nil
}

func (n *Navigator) GetParentPath() string {
	parentDir, _ := filepath.Split(n.currentPath)

	if parentDir == string(os.PathSeparator) {
		return parentDir
	}

	return strings.TrimSuffix(parentDir, string(os.PathSeparator))
}

func (n *Navigator) IncrementIndex() {
	if n.currentIndex < n.GetEntriesLength()-1 {
		n.currentIndex++
	}
}
