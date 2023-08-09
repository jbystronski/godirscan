package entry

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/jbystronski/godirscan/pkg/converter"
)

type Entry struct {
	Name  string
	IsDir bool
	Size  int
	Path  *string
}

func (e *Entry) FullPath() string {
	return filepath.Join(*e.Path, e.Name)
}

type Entries = []*Entry

func getTotalSize(path string, size *int) int {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if info.IsDir() {
		getTotalSize(filepath.Join(path, info.Name()), size)
	}

	*size += int(info.Size())

	return *size
}

func SetSort(seq *uint8, entries []*Entry) {
	if *seq == 0 {
		*seq = SortByName(entries)
	} else if *seq == 1 {
		*seq = SortBySizeAsc(entries)
	} else if *seq == 2 {
		*seq = sortBySizeDesc(entries)
	} else if *seq == 3 {
		*seq = SortByType(entries)
	}
}

func sortBySizeDesc(entries []*Entry) uint8 {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Size > entries[j].Size
	})

	return 3
}

func SortBySizeAsc(entries []*Entry) uint8 {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Size < entries[j].Size
	})

	return 2
}

func SortByName(entries []*Entry) uint8 {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir && !entries[j].IsDir {
			return true
		} else if !entries[i].IsDir && entries[j].IsDir {
			return false
		} else {
			return entries[i].Name < entries[j].Name
		}
	})

	return 1
}

func SortByType(entries []*Entry) uint8 {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir && !entries[j].IsDir {
			return true
		} else if !entries[i].IsDir && entries[j].IsDir {
			return false
		} else {
			return filepath.Ext(entries[i].Name) < filepath.Ext(entries[j].Name)
		}
	})
	return 0
}

func FormatSize(bytes int) string {
	if bytes < int(converter.KbInBytes) {
		return fmt.Sprintf("%d %s", bytes, converter.StorageUnits[0])
	}

	floatSize, unit := converter.BytesToFloat(bytes)

	return fmt.Sprintf("%.2f %s", floatSize, unit)
}

func (e Entry) PrintSize() string {
	return PrintSizeAsString(e.Size)
}

func PrintSizeAsString(size int) string {
	return fmt.Sprintf("%v", "["+FormatSize(size)+"]")
}

// func (e EntryChild) PrintPath(i int) string {
// 	fmtPath := ""

// 	if e.IsDir {
// 		fmtPath = fmtDir + e.Parent.Path + e.Name + terminal.ResetFmt
// 	} else {
// 		fmtPath = fmtDir + e.Parent.Path + fmtFile + e.Name + terminal.ResetFmt
// 	}

// 	return fmt.Sprintf("%s %v [%d]", fmtPath, formatSize(e.Size), i)
// }

// func (e Entry) String() string {
// 	indentParts := strings.Split(e.Pattern, "")
// 	format := ""
// 	for i, indent := range indentParts {
// 		if i == len(indentParts)-1 {
// 			if indent == "0" {
// 				format += terminal.CornerLine + terminal.Hseparator
// 			} else {
// 				format += terminal.TeeLine + terminal.Hseparator
// 			}
// 		} else {
// 			if indent == "0" {
// 				format += terminal.EmptyIndent
// 			} else {
// 				format += terminal.Vseparator + terminal.EmptyIndent
// 			}
// 		}
// 	}

// 	nameFmt := ""

// 	if e.IsDir {
