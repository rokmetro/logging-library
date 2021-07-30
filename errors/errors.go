package errors

import (
	"fmt"
	"strings"

	"github.com/rokmetro/logging-library/logutils"
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

//ErrorData generates an error for a data element
//	status: The status of the data
//	dataType: The data type that the error is occurring on
//	args: Any args that should be included in the message (nil if none)
func ErrorData(status logutils.MessageDataStatus, dataType logutils.MessageDataType, args logutils.MessageArgs) error {
	message := logutils.MessageData(status, dataType, args)
	message = strings.ToLower(message)
	return fmt.Errorf("%s() %s", getErrorPrevFuncName(), message)
}

//WrapErrorData wraps an error for a data element
//	status: The status of the data
//	dataType: The data type that the error is occurring on
//	args: Any args that should be included in the message (nil if none)
//  err: Error to wrap
func WrapErrorData(status logutils.MessageDataStatus, dataType logutils.MessageDataType, args logutils.MessageArgs, err error) error {
	message := logutils.MessageData(status, dataType, args)
	message = strings.ToLower(message)
	return fmt.Errorf("%s() %s", getErrorPrevFuncName(), message)
}

//ErrorAction generates an error for an action
//	action: The action that is occurring
//	dataType: The data type that the action is occurring on
//	args: Any args that should be included in the message (nil if none)
func ErrorAction(action logutils.MessageActionType, dataType logutils.MessageDataType, args logutils.MessageArgs) error {
	message := logutils.MessageAction(logutils.StatusError, action, dataType, args)
	message = strings.ToLower(message)
	return fmt.Errorf("%s() %s", getErrorPrevFuncName(), message)
}

//WrapErrorAction wraps an error for an action
//	action: The action that is occurring
//	dataType: The data type that the action is occurring on
//	args: Any args that should be included in the message (nil if none)
//	err: Error to wrap
func WrapErrorAction(action logutils.MessageActionType, dataType logutils.MessageDataType, args logutils.MessageArgs, err error) error {
	message := logutils.MessageAction(logutils.StatusError, action, dataType, args)
	message = strings.ToLower(message)
	return fmt.Errorf("%s() %s: %v", getErrorPrevFuncName(), message, err)
}

//getErrorPrevFuncName - fetches the previous function name for error functions
func getErrorPrevFuncName() string {
	return logutils.GetFuncName(4)
}
