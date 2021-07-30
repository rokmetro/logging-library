package logs

import (
	"strings"
)

type logLevel string

func LogLevelFromString(level string) *logLevel {
	var lLevel logLevel

	switch strings.ToLower(level) {
	case strings.ToLower(string(Debug)):
		lLevel = Debug
	case strings.ToLower(string(Info)):
		lLevel = Info
	case strings.ToLower(string(Warn)):
		lLevel = Warn
	case strings.ToLower(string(Error)):
		lLevel = Error
	}

	return &lLevel
}

const (
	//Levels
	Info  logLevel = "Info"
	Debug logLevel = "Debug"
	Warn  logLevel = "Warn"
	Error logLevel = "Error"
)
