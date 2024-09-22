package logger

import (
	"fmt"
	"io"
)

type Logger interface {
	Print(string)
	Info(string)
	Error(string)
}

type log struct {
	writer io.Writer
}

func (l log) Print(msg string) {

}

func (l log) Info(msg string) {
	fmt.Printf("Info: %s\n", msg)
}

func (l log) Error(msg string) {
	fmt.Printf("Error: %s\n", msg)
}

func New(writer io.Writer) Logger {
	return log{writer}
}
