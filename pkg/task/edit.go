package task

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Edit(fPath, editor string) {
	cmd := exec.Command(editor, fPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Peek(path string) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	r := bufio.NewReader(f)

	for {
		l, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(l)
	}
}
