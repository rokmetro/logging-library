package loglib

import (
	"fmt"
	"strings"
)

func WrapError(message string, err error) error {
	message = strings.ToLower(message)
	return fmt.Errorf("%s: %v", message, err)
}

func WrapErrorf(format string, err error, a ...interface{}) error {
	format = strings.ToLower(format)
	message := fmt.Sprintf(format, a)
	return fmt.Errorf("%s: %v", message, err)
}
