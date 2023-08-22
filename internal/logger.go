package internal

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	red    = 31
	green  = 32
	yellow = 33
	blue   = 36
	gray   = 37
)

const (
	PanicLevel logrus.Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type LoggerFormatter struct {
	TimeFormat string
}

// Format gather all the data from the log entry and format it to return the final log message
func (f *LoggerFormatter) Format(e *logrus.Entry) ([]byte, error) {
	msg := f.printColored(e)

	if len(e.Data) != 0 {
		var fields string
		for k, v := range e.Data {
			if fields == "" {
				fields = fmt.Sprintf("%s=%v", k, v)
			} else {
				fields = fmt.Sprintf("%s %s=%v", fields, k, v)
			}
		}

		msg = fmt.Sprintf("%s %s %s\n", msg, e.Message, fields)
	} else {
		msg = fmt.Sprintf("%s %s\n", msg, e.Message)
	}

	if f.TimeFormat != "" {
		msg = fmt.Sprintf("[%s] %s", time.Now().Format(f.TimeFormat), msg)
	}

	return []byte(msg), nil
}

// printColored returns the prefix for log message with correct log level color and name
func (f *LoggerFormatter) printColored(e *logrus.Entry) string {
	var levelColor int
	var levelText string
	switch e.Level {
	case TraceLevel:
		levelText = "TRACE:"
		levelColor = gray
	case DebugLevel:
		levelText = "DEBUG:"
		levelColor = blue
	case InfoLevel:
		levelText = "INFO:"
		levelColor = green
	case WarnLevel:
		levelText = "WARNING:"
		levelColor = yellow
	case ErrorLevel:
		levelText = "ERROR:"
		levelColor = red
	case FatalLevel:
		levelText = "FATAL:"
		levelColor = red
	case PanicLevel:
		levelText = "PANIC:"
		levelColor = red
	default:
		levelColor = blue
	}

	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", levelColor, levelText)
}
