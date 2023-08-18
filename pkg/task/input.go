package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	k "github.com/eiannone/keyboard"
	"github.com/jbystronski/godirscan/pkg/terminal"
)

func WaitInput(prompt, output string) (result string, err error) {
	print := func(s string) {
		terminal.ClearLine()
		terminal.CarriageReturn()
		fmt.Print(terminal.Prompt(prompt) + " " + s)
	}

	printChar := func(char, output string, cursorPos *int) string {
		output = output[0:*cursorPos] + strings.TrimSpace(char) + output[*cursorPos:]
		*cursorPos++

		return output
	}

	moveLeft := func(cursorPos *int) {
		if *cursorPos > 0 {
			terminal.MoveCursorLeft(1)
			*cursorPos--
		}
	}

	// promptLength := fmt.Sprint(len(prompt))
	cursorPosition := len(output)

	terminal.ClearLine()
	print(output)

	err = keyboard.Open()
	if err != nil {
		panic(err)
		return
	}

	defer func() {
		if keyboard.IsStarted(time.Millisecond * 10) {

			err := k.Close()
			if err != nil {
				panic(err)
			}
		}
	}()

	for {

		char, key, getKeyErr := k.GetKey()
		if getKeyErr != nil {
			err = getKeyErr
			return
			//	os.Exit(1)
		}

		switch key {
		case k.KeyEsc:
			terminal.ClearLine()
			terminal.CarriageReturn()
			return

		case k.KeyEnter:
			terminal.ClearLine()
			terminal.CarriageReturn()
			result = strings.TrimSpace(output)

			return

		case k.KeyArrowRight:

			if cursorPosition < len(output) {
				terminal.MoveCursorRight(1)
				cursorPosition++
			}

		case k.KeyArrowLeft:

			moveLeft(&cursorPosition)
		case k.KeyBackspace, k.KeyBackspace2:

			if len(output) > 0 {

				output = output[0:cursorPosition-1] + output[cursorPosition:]
				terminal.ClearLine()
				cursorPosition--

				print(output)
				terminal.MoveCursorLeft(len(output) - cursorPosition)
			}

		case k.KeyDelete:

			if len(output) > 0 && cursorPosition <= len(output)-1 {
				output = output[0:cursorPosition] + output[cursorPosition+1:]
				terminal.ClearLine()

				print(output)

				if len(output)-1-cursorPosition > 0 {
					terminal.MoveCursorLeft(len(output) - cursorPosition)
				} else {
					moveLeft(&cursorPosition)
				}

			}

		case k.KeySpace:

			output = printChar(" ", output, &cursorPosition)
			print(output)
			terminal.MoveCursorLeft(len(output) - cursorPosition)

		case k.KeyHome:

			print(output)
			cursorPosition = 0
			terminal.MoveCursorLeft(len(output))

		case k.KeyEnd:
			print(output)
			cursorPosition = len(output)
			terminal.MoveCursorRight(len(output) - cursorPosition)
			cursorPosition = len(output)

		default:

			if char != 0 {

				// c := terminal.RuneToUtf8String(char)

				output = printChar(string(char), output, &cursorPosition)
				print(output)

				terminal.MoveCursorLeft(len(output) - cursorPosition)
			}

		}

		// if key == k.KeyEsc {

		// 	terminal.ClearLine()
		// 	terminal.CarriageReturn()
		// 	return

		// } else if key == k.KeyArrowRight {
		// 	terminal.MoveCursorRight()
		// } else if key == k.KeyEnter {
		// 	terminal.ClearLine()
		// 	terminal.CarriageReturn()
		// 	output = strings.TrimSpace(output)
		// 	forwardOutput(output)
		// 	break

		// } else if key == k.KeyBackspace || key == k.KeyBackspace2 {
		// 	if len(output) > 0 {
		// 		output = output[:len(output)-1]
		// 		terminal.ClearLine()

		// 		print(output)
		// 	}
		// } else if key == k.KeySpace {
		// 	output += " "
		// 	print(output)
		// } else if char != 0 {
		// 	output += string(char)

		// 	print(output)
		// }

	}
}

// func listenQuit(quit chan<- struct{}) {
// 	for {
// 		key := make([]byte, 1)

// 		_, err := os.Stdin.Read(key)
// 		if err != nil {
// 			os.Exit(1)
// 		}

// 		if key[0] == 27 {
// 			fmt.Println("escape pressed")
// 			quit <- struct{}{}
// 			break
// 		}

// 	}
// }

// var wg sync.WaitGroup

// func WaitUserInput(prompt, placeholder string, forwardOutput func(string)) {
// 	// blockChan := make(chan struct{})

// 	// blockCchan <- struct{}{}

// 	// var wg sync.WaitGroup
// 	// wg.Add(2)

// 	// quit := make(chan struct{})

// 	// go func() {
// 	// 	defer wg.Done()
// 	// 	listenInput(prompt, output, quit)
// 	// }()

// 	// go func() {
// 	// 	defer wg.Done()
// 	// 	listenQuit(quit)
// 	// }()

// 	// wg.Wait()

// 	// forwardOutput("/home/kb")

// 	var input string

// 	// done := make(chan struct{}, 1)

// 	wg.Add(1)

// 	go func() {
// 		defer wg.Done()
// 		listenInput(prompt, placeholder, input)
// 	}()

// 	// wg.Add(1)
// 	// go func() {
// 	// 	defer wg.Done()
// 	// 	listenKeys(done)
// 	// }()

// 	wg.Wait()
// }

// // func listenKeys(done chan<- struct{}) {
// // 	defer keyboard.Close()

// // 	keyEvent, err := keyboard.GetKeys(1)
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	for {
// // 		select {
// // 		case event := <-keyEvent:
// // 			switch event.Key {
// // 			case keyboard.KeyEsc:
// // 				done <- struct{}{}
// // 				return
// // 			}
// // 		}
// // 	}
// // }

// func listenInput(prompt, defaultOption, input string) string {
// 	fmt.Print(prompt + defaultOption)

// 	fmt.Scanln(&input)

// 	return input
// }

// func listenKeys(done chan<- struct{}) {
// 	for {
// 		key := make([]byte, 1)

// 		_, err := os.Stdin.Read(key)
// 		if err != nil {
// 			os.Exit(1)
// 		}

// 		if key[0] == 27 {
// 			key[0] = 13

// 			fmt.Println("escape pressed")
// 			done <- struct{}{}
// 			return
// 		}

// 	}
// }

// // func listenInput(prompt, placeholder string, done <-chan struct{}) {
// // 	_, err := os.Stdin.Write([]byte(placeholder))
// // 	if err != nil {
// // 		os.Exit(1)
// // 	}

// // 	reader := bufio.NewReader(os.Stdin)

// // 	for {
// // 		select {
// // 		case <-done:
// // 			return
// // 		default:
// // 			input, err := reader.ReadString('\n')
// // 			if err != nil {
// // 				os.Exit(1)
// // 			}
// // 			fmt.Printf("%s%s", prompt, input)
// // 		}
// // 	}
// // }
