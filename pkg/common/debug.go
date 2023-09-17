package common

import (
	"bufio"
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
	k "github.com/eiannone/keyboard"
)

func Log(msg ...interface{}) {
	msgToString := fmt.Sprintln(msg...)

	file, err := os.OpenFile("/home/kb/log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(msgToString)
	if err != nil {
		panic(err)
	}
}

func ViewLog() {
	ClearScreen()
	filePath := "/home/kb/log"

	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a bufio.Scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate through the lines
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line) // Process each line as needed
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	for {

		_, key, getKeyErr := k.GetKey()
		if getKeyErr != nil {
			panic(getKeyErr)
		}

		switch key {
		case keyboard.KeyEsc:
			return
		}

	}
}
