package task

import (
	"os"
	"os/exec"

	"github.com/jbystronski/godirscan/pkg/config"
)

func Edit(fPath string) error {
	cmd := exec.Command(config.Cfg.DefaultEditor, fPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
