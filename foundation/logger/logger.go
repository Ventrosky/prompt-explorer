package logger

import (
	"log"
	"os"
)

// Logger provides logging functionality
type Logger struct {
	info  *log.Logger
	error *log.Logger
}

// New creates a new Logger instance
func New() *Logger {
	return &Logger{
		info:  log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lshortfile),
		error: log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile),
	}
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.info.Println(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.error.Println(msg)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.info.Printf(format, args...)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.error.Printf(format, args...)
}
