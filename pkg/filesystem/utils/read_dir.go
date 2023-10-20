package utils

import (
	"io/fs"
	"os"
	"os/user"
	"runtime"
	"strings"
)

func ReadIgnorePermission(path string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsPermission(err) {
			return []fs.DirEntry{}, nil
		} else {
			return nil, err
		}
	}

	return entries, nil
}

func ResolveUserDirectory(path *string) {
	switch runtime.GOOS {
	case "darwin":
		{
			break
		}
	case "windows":
		{
			break
		}
	default:
		{
			if strings.HasPrefix(*path, "~") {
				currentUser, err := user.Current()
				if err != nil {
					panic(err)
				}

				*path = strings.Replace(*path, "~", currentUser.HomeDir, 1)

			}
		}
	}
}
