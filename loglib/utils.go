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

//MessageData generates a message string for a data element
//	status: The status of the data
//	dataType: The data type
//	args: Any args that should be included in the message (nil if none)
func MessageData(status logDataStatus, dataType LogData, args logArgs) string {
	argStr := ""
	if args != nil {
		argStr = args.String()
		if argStr != "" {
			argStr = ": " + argStr
		}
	}

	return fmt.Sprintf("%s %s%s", status, dataType, argStr)
}

//ErrorData generates an error for a data element
//	status: The status of the data
//	dataType: The data type that the error is occurring on
//	args: Any args that should be included in the message (nil if none)
func ErrorData(status logDataStatus, dataType LogData, args logArgs) error {
	message := MessageData(status, dataType, args)
	message = strings.ToLower(message)
	return fmt.Errorf("%s() %s", getErrorPrevFuncName(), message)
}

//WrapErrorData wraps an error for a data element
//	status: The status of the data
//	dataType: The data type that the error is occurring on
//	args: Any args that should be included in the message (nil if none)
//  err: Error to wrap
func WrapErrorData(status logDataStatus, dataType LogData, args logArgs, err error) error {
	message := MessageData(status, dataType, args)
	message = strings.ToLower(message)
	return fmt.Errorf("%s() %s", getErrorPrevFuncName(), message)
}

//MessageAction generates a message string for an action
//	status: The status of the action
//	action: The action that is occurring
//	dataType: The data type that the action is occurring on
//	args: Any args that should be included in the message (nil if none)
func MessageAction(status logActionStatus, action LogAction, dataType LogData, args logArgs) string {
	argStr := ""
	if args != nil {
		argStr = args.String()
		if argStr != "" {
			argStr = " for " + argStr
		}
	}

	return fmt.Sprintf("%s %s %s%s", status, action, dataType, argStr)
}

//ErrorAction generates an error for an action
//	action: The action that is occurring
//	dataType: The data type that the action is occurring on
//	args: Any args that should be included in the message (nil if none)
func ErrorAction(action LogAction, dataType LogData, args logArgs) error {
	message := MessageAction(StatusError, action, dataType, args)
	message = strings.ToLower(message)
	return fmt.Errorf("%s() %s", getErrorPrevFuncName(), message)
}

//WrapErrorAction wraps an error for an action
//	action: The action that is occurring
//	dataType: The data type that the action is occurring on
//	args: Any args that should be included in the message (nil if none)
//	err: Error to wrap
func WrapErrorAction(action LogAction, dataType LogData, args logArgs, err error) error {
	message := MessageAction(StatusError, action, dataType, args)
	message = strings.ToLower(message)
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

//getLogPrevFuncName - fetches the calling function name when logging
//	layer: Number of internal library function calls above caller
func getLogPrevFuncName(layer int) string {
	return GetFuncName(5 + layer)
}

//GetFuncName fetches the name of a function caller based on the numFrames
func GetFuncName(numFrames int) string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(numFrames, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}
