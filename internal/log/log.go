package log

import (
	"io"
	"log"
)

var (
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger

	Debug   func(format string, args ...interface{})
	Info    func(format string, args ...interface{})
	Warning func(format string, args ...interface{})
	Error   func(format string, args ...interface{})
)

func init() { Init(nil) }

func Init(output io.Writer) {
	writer := log.Default().Writer()
	if output != nil {
		writer = output
	}

	debugLogger = log.New(writer, "DEBUG: ", log.LstdFlags)
	Debug = debugLogger.Printf

	infoLogger = log.New(writer, "INFO: ", log.LstdFlags)
	Info = infoLogger.Printf

	warningLogger = log.New(writer, "WARNING: ", log.LstdFlags)
	Warning = warningLogger.Printf

	errorLogger = log.New(writer, "ERROR: ", log.LstdFlags|log.Lshortfile)
	Error = errorLogger.Printf
}
