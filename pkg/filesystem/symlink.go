package filesystem

import "fmt"

type FsSymlink struct {
	FsEntry
	FsFile
}

func (sym *FsSymlink) String() string {
	return fmt.Sprint(sym.Name(), " ", "[sym]")
}

func (sym *FsSymlink) execute() error {
	fmt.Print("following symlink")
	return nil
}
