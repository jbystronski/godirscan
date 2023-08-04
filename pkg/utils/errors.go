package utils

import (
	"fmt"
	"os"
)

func printDefaultErrorAndExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
