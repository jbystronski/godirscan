package filesystem

import "fmt"

type FsDirectory struct {
	FsEntry
}

func (f *FsDirectory) String() string {
	return fmt.Sprint(f.Name())
}
