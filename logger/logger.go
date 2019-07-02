package logger

import (
	"fmt"
	"io"
	"log"
)

var (
	info    *log.Logger
	warning *log.Logger
	err     *log.Logger
)

func Init(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime)

	warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime)

	err = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime)
}

func Info(msg string) {
	info.Println(msg)
}

func Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v)
	warning.Println(msg)
}

func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v)
	err.Println(msg)
}
