package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	level  LogLevel
	logger *log.Logger
}

var std *Logger

// Init initializes the global logger
func Init(levelStr string) {
	level := parseLogLevel(levelStr)
	std = &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// parseLogLevel converts string to LogLevel
func parseLogLevel(levelStr string) LogLevel {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	if std != nil && std.level <= DEBUG {
		std.logger.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	if std != nil && std.level <= INFO {
		std.logger.Printf("[INFO] "+format, v...)
	}
}

// Warn logs a warning message
func Warn(format string, v ...interface{}) {
	if std != nil && std.level <= WARN {
		std.logger.Printf("[WARN] "+format, v...)
	}
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	if std != nil && std.level <= ERROR {
		std.logger.Printf("[ERROR] "+format, v...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(format string, v ...interface{}) {
	if std != nil {
		std.logger.Fatalf("[FATAL] "+format, v...)
	} else {
		log.Fatalf("[FATAL] "+format, v...)
	}
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if std == nil {
		Init("INFO")
	}
	return std
}

// WithPrefix returns a new logger with a prefix
func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{
		level:  l.level,
		logger: log.New(os.Stdout, fmt.Sprintf("[%s] ", prefix), log.LstdFlags),
	}
}
