package logging

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	red    = 31
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

type LoggerFormatter struct{}

// Format gather all the data from the log entry and format it to return the final log message
func (f *LoggerFormatter) Format(e *logrus.Entry) ([]byte, error) {
	msg := f.printColored(e)

	if len(e.Data) != 0 {
		for k, v := range e.Data {
			msg = fmt.Sprintf("%s %s=%v", msg, k, v)
		}

		msg = fmt.Sprintf("%s message='%s'\n", msg, e.Message)
	} else {
		msg = fmt.Sprintf("%s %s\n", msg, e.Message)
	}

	return []byte(msg), nil
}

// printColored returns the prefix for log message with correct log level color and name
func (f *LoggerFormatter) printColored(e *logrus.Entry) string {
	var levelColor int
	var levelText string
	switch e.Level {
	case TraceLevel:
		levelText = "trace:"
		levelColor = gray
	case DebugLevel:
		levelText = "debug:"
		levelColor = gray
	case InfoLevel:
		levelText = "info:"
		levelColor = blue
	case WarnLevel:
		levelText = "warn:"
		levelColor = yellow
	case ErrorLevel:
		levelText = "error:"
		levelColor = red
	case FatalLevel:
		levelText = "fatal:"
		levelColor = red
	case PanicLevel:
		levelText = "panic:"
		levelColor = red
	default:
		levelColor = blue
	}

	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", levelColor, levelText)
}
