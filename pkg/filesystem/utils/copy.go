package utils

import (
	"io"
	"os"
)

func Copy(sourcePath, targetPath string) error {
	sourceFile, err := os.Open(sourcePath)
	defer sourceFile.Close()

	if err != nil {
		return err
	}

	targetFile, err := os.Create(targetPath)

	defer targetFile.Close()

	if err != nil {
		return err
	}

	_, err = io.Copy(targetFile, sourceFile)

	if err != nil {
		return err
	}

	info, _ := os.Stat(sourcePath)

	err = os.Chmod(targetPath, info.Mode())
	if err != nil {
		return err
	}
	return nil
}
