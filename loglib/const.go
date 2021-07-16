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

func (s StringArgs) String() string {
	return string(s)
}

type logLevel string

type logDataStatus string
type logActionStatus string

type LogAction string

type LogData string

const (
	//Levels
	Info  logLevel = "Info"
	Debug logLevel = "Debug"
	Warn  logLevel = "Warn"
	Error logLevel = "Error"

	//Errors
	Unimplemented string = "Unimplemented"

	//Types
	ValidStatus   logDataStatus = "Valid"
	FoundStatus   logDataStatus = "Found"
	InvalidStatus logDataStatus = "Invalid"
	MissingStatus logDataStatus = "Missing"

	SuccessStatus logActionStatus = "Success"
	ErrorStatus   logActionStatus = "Error"

	//Data
	ArgData LogData = "arg"

	//Primitives
	IntData    LogData = "int"
	UintData   LogData = "uint"
	FloatData  LogData = "float"
	BoolData   LogData = "bool"
	StringData LogData = "string"
	ByteData   LogData = "byte"
	ErrorData  LogData = "error"

	//Requests
	RequestBodyData  LogData = "request body"
	ResponseData     LogData = "response"
	ResponseBodyData LogData = "response body"
	QueryParamData   LogData = "query param"

	//Auth
	TokenData      LogData = "token"
	ClaimsData     LogData = "claims"
	ClaimData      LogData = "claim"
	ScopeData      LogData = "scope"
	PermissionData LogData = "permission"

	//Actions
	InitializeAction LogAction = "initializing"
	ComputeAction    LogAction = "computing"
	RegisterAction   LogAction = "registering"

	//Encryption
	EncryptAction LogAction = "entrypting"
	DecryptAction LogAction = "decrypting"

	//Request/Response Actions
	SendAction LogAction = "sending"
	ReadAction LogAction = "reading"

	//Encode Actions
	ParseAction  LogAction = "parsing"
	EncodeAction LogAction = "encoding"
	DecodeAction LogAction = "decoding"

	//Marshal Actions
	MarshalAction   LogAction = "marshalling"
	UnmarshalAction LogAction = "unmarshalling"
	ValidateAction  LogAction = "validating"
	CastAction      LogAction = "casting to"

	//Cache Actions
	CacheAction     LogAction = "caching"
	LoadCacheAction LogAction = "loading cached"

	//Operation Actions
	GetAction    LogAction = "getting"
	CreateAction LogAction = "creating"
	UpdateAction LogAction = "updating"
	DeleteAction LogAction = "deleting"

	//Storage Actions
	FindAction    LogAction = "finding"
	InsertAction  LogAction = "inserting"
	ReplaceAction LogAction = "replacing"
	SaveAction    LogAction = "saving"
	CountAction   LogAction = "counting"
)
