package logger

import (
	"log"
	"os"
)

type Logger struct {
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Debug *log.Logger
}

var logInstance *Logger

// InitLogger initializes the logger with the specified level and sets the logInstance
func InitLogger(level string) *Logger {
	out := os.Stdout

	debug := log.New(out, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	info := log.New(out, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warn := log.New(out, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	err := log.New(out, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Open /dev/null to discard unwanted log output
	devNull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0644)

	logger := &Logger{
		Debug: debug,
		Info:  info,
		Warn:  warn,
		Error: err,
	}

	// Configure the logger based on the log level
	switch level {
	case "debug":
		// Debug level includes all logs (no need to change anything)
	case "info":
		// Info level excludes debug logs
		logger.Debug.SetOutput(devNull) // Disable debug logs
	case "warn":
		// Warn level excludes both debug and info logs
		logger.Debug.SetOutput(devNull)
		logger.Info.SetOutput(devNull)
	case "error":
		// Error level excludes all logs except errors
		logger.Debug.SetOutput(devNull)
		logger.Info.SetOutput(devNull)
		logger.Warn.SetOutput(devNull)
	default:
		// Default to info level if no valid level is provided
		logger.Debug.SetOutput(devNull)
	}

	// Assign the logger to the global logInstance variable
	logInstance = logger

	return logger
}

// GetLogger returns the initialized logger
func GetLogger(level string) *Logger {
	if logInstance == nil {
		InitLogger(level)
	}
	return logInstance
}
