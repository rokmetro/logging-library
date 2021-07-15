package loglib

import (
	"errors"
	"fmt"
	"strings"
)

//NewError returns an error containing the provided message
func NewError(message string) error {
	message = strings.ToLower(message)
	return errors.New(message)
}

//FmtError returns an error containing the formatted message
func FmtError(message string, a ...interface{}) error {
	message = strings.ToLower(message)
	return fmt.Errorf(message, a)
}

//WrapErrorf returns an error containing the provided message and error
func WrapError(message string, err error) error {
	message = strings.ToLower(message)
	return fmt.Errorf("%s: %v", message, err)
}

//WrapErrorf returns an error containing the formatted message and provided error
func WrapErrorf(format string, err error, a ...interface{}) error {
	format = strings.ToLower(format)
	message := fmt.Sprintf(format, a)
	return fmt.Errorf("%s: %v", message, err)
}

func containsString(slice []string, val string) bool {
	for _, v := range slice {
		if val == v {
			return true
		}
	}
	return false
}
