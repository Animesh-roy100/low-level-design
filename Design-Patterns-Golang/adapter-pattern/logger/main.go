package logger

import "fmt"

// Real Use: Legacy system integration, third-party API wrappers
// Key Usage: Interface compatibility

type OldLogger interface { // Legacy interface
	LogMessage(msg string)
}

type NewLogger interface { // Modern interface
	Log(level string, message string)
}

// Modern implementation
type CloudLogger struct{}

func (c *CloudLogger) Log(level, msg string) {
	fmt.Printf("[%s] %s\n", level, msg)
}

type LoggerAdapter struct {
	CloudLogger *CloudLogger
}

func (a *LoggerAdapter) LogMessage(msg string) {
	a.CloudLogger.Log("INFO", msg) // Adapt old call
}

// Usage:
// oldLogger := &LoggerAdapter{CloudLogger: &CloudLogger{}}
// oldLogger.LogMessage("Legacy system message")
