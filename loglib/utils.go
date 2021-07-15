package loglib

import (
	"fmt"
	"runtime"
	"strings"
)

//NewError returns an error containing the provided message
func NewError(message string) error {
	message = strings.ToLower(message)
	return fmt.Errorf("%s() %s", getErrorPrevFuncName(), message)
}

//NewErrorf returns an error containing the formatted message
func NewErrorf(message string, args ...interface{}) error {
	message = strings.ToLower(message)
	message = fmt.Sprintf(message, args...)
	return fmt.Errorf("%s() %s", getErrorPrevFuncName(), message)
}

//WrapErrorf returns an error containing the provided message and error
func WrapError(message string, err error) error {
	message = strings.ToLower(message)
	return fmt.Errorf("%s() %s: %v", getErrorPrevFuncName(), message, err)
}

//WrapErrorf returns an error containing the formatted message and provided error
func WrapErrorf(format string, err error, args ...interface{}) error {
	format = strings.ToLower(format)
	message := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s() %s: %v", getErrorPrevFuncName(), message, err)
}

func containsString(slice []string, val string) bool {
	for _, v := range slice {
		if val == v {
			return true
		}
	}
	return false
}

//getErrorPrevFuncName - fetches the previous function name for error functions
func getErrorPrevFuncName() string {
	return GetFuncName(4)
}

//getRequestFieldsPrevFuncName - fetches the previous function name for the GetRequestFields function
func getRequestFieldsPrevFuncName() string {
	return GetFuncName(5)
}

//GetFuncName fetches the name of a function caller based on the numFrames
func GetFuncName(numFrames int) string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(numFrames, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}
