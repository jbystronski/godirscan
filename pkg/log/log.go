package log

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

type Logger struct {
	mu   *sync.RWMutex
	path string
	ch   chan (string)
}

func New(path string, truncate bool) (error, *Logger) {
	_, err := os.Stat(path)

	if err != nil && errors.Is(err, os.ErrNotExist) {

		_, err := os.Create(path)
		if err != nil {
			return err, nil
		}

	}

	if truncate {
		err := os.Truncate(path, 0)
		if err != nil {
			return err, nil
		}
	}

	return nil, &Logger{mu: &sync.RWMutex{}, path: path}
}

func (l *Logger) Log(msg ...interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	file, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintln(msg...))
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) Clear() error {
	err := os.Truncate(l.path, 0)
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) Delete() error {
	err := os.RemoveAll(l.path)
	if err != nil {
		return err
	}

	return nil
}
