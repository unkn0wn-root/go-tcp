package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
}

func NewLogger() *Logger {
	flags := log.Ldate | log.Ltime | log.Lmicroseconds

	return &Logger{
		InfoLogger:  log.New(os.Stdout, "INFO: ", flags),
		ErrorLogger: log.New(os.Stderr, "ERROR: ", flags),
		DebugLogger: log.New(os.Stdout, "DEBUG: ", flags),
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.InfoLogger.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.ErrorLogger.Printf(format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		l.DebugLogger.Printf(format, v...)
	}
}

func (l *Logger) LogToFile(filename string) error {
	file, err := os.OpenFile(
		filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	l.InfoLogger.SetOutput(file)
	l.ErrorLogger.SetOutput(file)
	l.DebugLogger.SetOutput(file)

	return nil
}
