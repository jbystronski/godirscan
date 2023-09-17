package filesystem

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/common"
)

type FsFile struct {
	FsEntry
}

func (f *FsFile) String() string {
	return fmt.Sprint(f.Name())
}

func (f *FsFile) execute() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		{
			cmd = exec.Command("open", f.FullPath())
		}
	case "windows":
		{
			cmd = exec.Command("cmd", "/c", "start", f.FullPath())
		}
	default:
		{
			cmd = exec.Command("xdg-open", f.FullPath())
		}
	}
	common.ClearScreen()
	fmt.Println(common.CurrentTheme.Accent + "Press esc to return, command execution output: " + "\033[0m")
	fmt.Println()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	for {

		_, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		if key == keyboard.KeyEsc {
			return nil
		}

	}
}

func (f *FsFile) Edit() error {
	cmd := exec.Command(common.Cfg.DefaultEditor, f.FullPath())

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
