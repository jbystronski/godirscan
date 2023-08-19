package terminal

import (
	"fmt"
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
