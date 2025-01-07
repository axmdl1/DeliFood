package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func NewLogger() *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{}) // Use JSON format for logs
	log.SetOutput(os.Stdout)                  // Output logs to standard output
	log.SetLevel(logrus.InfoLevel)            // Set default log level to INFO

	return &Logger{
		log: log,
	}
}

// Info logs an informational message
func (l *Logger) Info(message string, fields map[string]interface{}) {
	l.log.WithFields(logrus.Fields(fields)).Info(message)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields map[string]interface{}) {
	l.log.WithFields(logrus.Fields(fields)).Warn(message)
}

// Error logs an error message
func (l *Logger) Error(message string, fields map[string]interface{}) {
	l.log.WithFields(logrus.Fields(fields)).Error(message)
}

// Example usage:
// logger := logger.NewLogger()
// logger.Info("Application started", map[string]interface{}{"module": "main"})
