package logwrapper

import (
	"github.com/sirupsen/logrus"
)

type StandardLogger struct {
	*logrus.Logger
}

func NewLogger() *StandardLogger {
	var baseLogger = logrus.New()
	var standardLogger = &StandardLogger{baseLogger}
	standardLogger.Formatter = &logrus.JSONFormatter{}
	return standardLogger
}

// InvalidArgValue is a standard error message for invalid arguments
func (l *StandardLogger) InvalidArg(argumentName string, argumentValue string) {
	l.Errorf("Invalid value for argument", argumentName, argumentValue)
}

// MissingArg is a standard error message for missing arguments
func (l *StandardLogger) MissingArg(argumentName string) {
	l.Errorf("Missing argument", argumentName)
}
