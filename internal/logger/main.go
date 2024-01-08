package logger

import (
	"fmt"
	"strings"
)

// ClassicLogFormat is the format used to print user logs.
//
// Format: "LEVEL MESSAGE"
const ClassicLogFormat = "%s %s"

type Level uint8

const (
	Error Level = iota
	Warn
	Info
	Debug
)

var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

func (l Level) String() string {
	switch l {
	case Error:
		return fmt.Sprintf("%s%s%s", Red, "ERROR", Reset)
	case Warn:
		return fmt.Sprintf("%s%s%s", Yellow, "WARN", Reset)
	case Info:
		return fmt.Sprintf("%s%s%s", Green, "INFO", Reset)
	case Debug:
		return fmt.Sprintf("%s%s%s", Blue, "DEBUG", Reset)
	default:
		panic("invalid logger level")
	}
}

func LevelFromString(lvl string) Level {
	switch strings.ToUpper(lvl) {
	case "ERROR":
		return Error
	case "WARN":
		return Warn
	case "INFO":
		return Info
	case "DEBUG":
		return Debug
	default:
		panic("invalid logger level")
	}
}

type Logger struct {
	Lvl Level
}

// New instanciates a new Logger. It takes a Level as argument.
func New(level Level) *Logger {
	return &Logger{
		Lvl: level,
	}
}

// Error logs an error message.
func (l Logger) Error(msg string, args ...interface{}) {
	l.log(Error, msg, args...)
}

// Warn logs a warning message.
func (l Logger) Warn(msg string, args ...interface{}) {
	l.log(Warn, msg, args...)
}

// Info logs an info message.
func (l Logger) Info(msg string, args ...interface{}) {
	l.log(Info, msg, args...)
}

// Debug logs a debug message.
func (l Logger) Debug(msg string, args ...interface{}) {
	l.log(Debug, msg, args...)
}

func (l Logger) log(level Level, msg string, args ...interface{}) {
	if l.Lvl >= level {
		fmt.Printf(ClassicLogFormat, l.Lvl, fmt.Sprintf(msg, args...))
	}
}
