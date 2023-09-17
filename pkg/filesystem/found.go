package filesystem

import "fmt"

type FsFound struct {
	FsEntry
}

func (f *FsFound) String() string {
	return fmt.Sprint(f.FullPath())
}
