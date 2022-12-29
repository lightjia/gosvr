package core

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Log levels to control the logging output.
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

type Logger struct {
	level  int
	logger *log.Logger
}

var LOG = Logger{LevelTrace, log.New(os.Stdout, "", log.Ldate|log.Ltime)}

type brush func(string) string

func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []brush{
	newBrush("1;41"), // 红色底
	newBrush("1;35"), // 紫色
	newBrush("1;34"), // 蓝色
	newBrush("1;31"), // 红色
	newBrush("1;33"), // 黄色
	newBrush("1;36"), // 天蓝色
	newBrush("1;32"), // 绿色
}

func NewLogger(out io.Writer, prefix string, flag int, level int) *Logger {
	levelStr := os.Getenv("LIGHT_LOG_LEVEL")
	if len(levelStr) != 0 {
		if vl, err := strconv.Atoi(levelStr); err == nil {
			level = vl
		}
	}

	return &Logger{level: level, logger: log.New(os.Stdout, prefix, flag)}
}

// LogLevel returns the global log level and can be used in own implementations of the logger interface.
func (l *Logger) Level() int {
	return l.level
}

// SetLogLevel sets the global log level used by the simple logger.
func (l *Logger) SetLevel(lv int) {
	l.level = lv
}

// Trace logs a message at trace level.
func (l *Logger) Trace(v ...interface{}) {
	if l.level <= LevelTrace {
		l.logger.Printf("[Trace] %v", v)
	}
}

// Debug logs a message at debug level.
func (l *Logger) Debug(v ...interface{}) {
	if l.level <= LevelDebug {
		l.logger.Printf("%v", colors[6](fmt.Sprintf("[Debug] %v", v)))
	}
}

// Info logs a message at info level.
func (l *Logger) Info(v ...interface{}) {
	if l.level <= LevelInfo {
		l.logger.Printf("%v", colors[2](fmt.Sprintf("[Info] %v", v)))
	}
}

// Warning logs a message at warning level.
func (l *Logger) Warn(v ...interface{}) {
	if l.level <= LevelWarning {
		l.logger.Printf("%v", colors[4](fmt.Sprintf("[Warn] %v", v)))
	}
}

// Error logs a message at error level.
func (l *Logger) Error(v ...interface{}) {
	if l.level <= LevelError {
		l.logger.Printf("%v", colors[0](fmt.Sprintf("[Error] %v", v)))
	}
}

// Critical logs a message at critical level.
func (l *Logger) Critical(v ...interface{}) {
	if l.level <= LevelCritical {
		l.logger.Printf("%v", colors[0](fmt.Sprintf("[Critical] %v", v)))
	}
}
