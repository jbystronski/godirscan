package terminal

import (
	"fmt"
	"os"
	"time"
)

func FlashError(err error) {
	go func() {
		Cell(totalLines, 1)
		fmt.Print(err)
		time.Sleep(time.Second * 2)
		Cell(totalLines, 1)
		ClearLine()
		CarriageReturn()

		printHelpers()
		return
	}()
}

func HandleError(err error) {
	switch err {

	case os.ErrPermission:

	case os.ErrNotExist:

	case os.ErrExist:

	default:

	}
}
