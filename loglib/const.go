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
	StatusValid   logDataStatus = "Valid"
	StatusInvalid logDataStatus = "Invalid"
	StatusFound   logDataStatus = "Found"
	StatusMissing logDataStatus = "Missing"

	StatusSuccess logActionStatus = "Success"
	StatusError   logActionStatus = "Error"

	//Data
	TypeArg         LogData = "arg"
	TypeTransaction LogData = "transaction"

	//Primitives
	TypeInt    LogData = "int"
	TypeUint   LogData = "uint"
	TypeFloat  LogData = "float"
	TypeBool   LogData = "bool"
	TypeString LogData = "string"
	TypeByte   LogData = "byte"
	TypeError  LogData = "error"

	//Requests
	TypeRequest      LogData = "request"
	TypeRequestBody  LogData = "request body"
	TypeResponse     LogData = "response"
	TypeResponseBody LogData = "response body"
	TypeQueryParam   LogData = "query param"

	//Auth
	TypeToken      LogData = "token"
	TypeClaims     LogData = "claims"
	TypeClaim      LogData = "claim"
	TypeScope      LogData = "scope"
	TypePermission LogData = "permission"

	//Actions
	ActionInitialize LogAction = "initializing"
	ActionCompute    LogAction = "computing"
	ActionRegister   LogAction = "registering"
	ActionStart      LogAction = "starting"
	ActionCommit     LogAction = "committing"

	//Encryption Actions
	ActionEncrypt LogAction = "encrypting"
	ActionDecrypt LogAction = "decrypting"

	//Request/Response Actions
	ActionSend LogAction = "sending"
	ActionRead LogAction = "reading"

	//Encode Actions
	ActionParse  LogAction = "parsing"
	ActionEncode LogAction = "encoding"
	ActionDecode LogAction = "decoding"

	//Marshal Actions
	ActionMarshal   LogAction = "marshalling"
	ActionUnmarshal LogAction = "unmarshalling"
	ActionValidate  LogAction = "validating"
	ActionCast      LogAction = "casting to"

	//Cache Actions
	ActionCache     LogAction = "caching"
	ActionLoadCache LogAction = "loading cached"

	//Operation Actions
	ActionGet    LogAction = "getting"
	ActionCreate LogAction = "creating"
	ActionUpdate LogAction = "updating"
	ActionDelete LogAction = "deleting"

	//Storage Actions
	ActionFind    LogAction = "finding"
	ActionInsert  LogAction = "inserting"
	ActionReplace LogAction = "replacing"
	ActionSave    LogAction = "saving"
	ActionCount   LogAction = "counting"
)
