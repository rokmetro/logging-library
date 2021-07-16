package loglib

import (
	"fmt"
	"strings"
)

type logArgs interface {
	String() string
}

type FieldArgs Fields

func (f *FieldArgs) String() string {
	if f == nil {
		return ""
	}

	argMsg := ""
	for k, v := range *f {
		if argMsg != "" {
			argMsg += ", "
		}

		if v != nil {
			argMsg += fmt.Sprintf("%s=%v", k, v)
		} else {
			argMsg += k
		}
	}

	return argMsg
}

type ListArgs []string

func (l *ListArgs) String() string {
	if l == nil {
		return ""
	}

	return strings.Join(*l, ", ")
}

type StringArgs string

func (s *StringArgs) String() string {
	if s == nil {
		return ""
	}

	return string(*s)
}

type logLevel string

type logDataStatus string
type logActionStatus string

type logAction string

//NewErrorAction creates a new errorAction type from the provided string
func NewErrorAction(action string) logAction {
	return logAction(action)
}

type logData string

//NewErrorData creates a new errorData type from the provided string
func NewErrorData(dataType string) logData {
	return logData(dataType)
}

const (
	//Levels
	Info  logLevel = "Info"
	Debug logLevel = "Debug"
	Warn  logLevel = "Warn"
	Error logLevel = "Error"

	//Types
	ValidStatus   logDataStatus = "Valid"
	FoundStatus   logDataStatus = "Found"
	InvalidStatus logDataStatus = "Invalid"
	MissingStatus logDataStatus = "Missing"

	SuccessStatus logActionStatus = "Success"
	ErrorStatus   logActionStatus = "Error"

	//Data
	RequestData      logData = "request"
	RequestBodyData  logData = "request body"
	ResponseData     logData = "response"
	ResponseBodyData logData = "response body"
	QueryParamData   logData = "query param"
	ArgData          logData = "arg"

	//Request/Response Actions
	MakeAction logAction = "making"
	ReadAction logAction = "reading"

	//Encode Actions
	MarshalAction   logAction = "marshalling"
	UnmarshalAction logAction = "unmarshalling"
	ValidateAction  logAction = "validating"
	CastAction      logAction = "casting to"

	//Operation Actions
	GetAction    logAction = "getting"
	CreateAction logAction = "creating"
	UpdateAction logAction = "updating"
	DeleteAction logAction = "deleting"

	//Storage Actions
	FindAction    logAction = "finding"
	InsertAction  logAction = "inserting"
	ReplaceAction logAction = "replacing"
	SaveAction    logAction = "saving"
	CountAction   logAction = "counting"
)
