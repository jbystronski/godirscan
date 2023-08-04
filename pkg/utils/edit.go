package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func edit(fPath, editor string, errHandle func(err error)) {
	cmd := exec.Command(editor, fPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	errHandle(err)
}

func peek(path string) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)

	printDefaultErrorAndExit(err)

	defer f.Close()

	r := bufio.NewReader(f)

	for {
		l, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}

		printDefaultErrorAndExit(err)

		fmt.Println(l)
	}
}
