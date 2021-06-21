package logging

import (
	"runtime"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type StandardLogger struct {
	*logrus.Entry
}

type Log struct {
	*StandardLogger
	traceID    string
	spanID     string
	prevSpanID string
	context    map[string]interface{}
}

//Constructor for new logger object with configuration at the service level
func NewLogger(serviceName string) *StandardLogger {
	var baseLogger = logrus.New()
	baseLogger.Formatter = &logrus.JSONFormatter{}
	standardFields := logrus.Fields{"serviceName": serviceName} //All common fields for logs of a given service
	contextLogger := &StandardLogger{baseLogger.WithFields(standardFields)}
	return contextLogger
}

//Constructor for a new log object of a request
func NewLog(logger *StandardLogger, traceID string, prevSpanID string, stackTrace []string, context map[string]interface{}) *Log {
	if traceID == "" {
		traceID = uuid.New().String()
	}
	spanID := uuid.New().String()
	l := &Log{logger, traceID, spanID, prevSpanID, context}
	return l
}

//getRequestFields() populates a map with all the fields of a request
func (l *Log) getRequestFields() logrus.Fields {
	fields := logrus.Fields{"trace_id": l.traceID, "span_id": l.spanID,
		"prev_span_id": l.prevSpanID, "function_name": getPrevFuncName(), "context": l.context}
	return fields
}

// InvalidArgValue is a standard error interface for invalid arguments
func (l *Log) InvalidArg(argumentName string, argumentValue interface{}) {
	fields := l.getRequestFields()
	fields["argument"] = argumentName
	fields["value"] = argumentValue
	l.WithFields(fields).Error("Invalid argument")
}

// MissingArg is a standard error interface for missing arguments
func (l *Log) MissingArg(argumentName string) {
	fields := l.getRequestFields()
	fields["argument"] = argumentName
	l.WithFields(fields).Error("Missing argument")
}

//TODO: More error interfaces to be added

func (l *Log) ErrorWithFields(message string, internal string) {
	fields := logrus.Fields{}
	l.WithFields(fields).Error(message)
}

//AddContext adds any relevant unstructured data to context map
func (l *Log) AddContext(fieldName string, value interface{}) {
	l.context[fieldName] = value
}

//PrintContext prints the entire context of a log object
func (l *Log) PrintContext() {
	fields := l.getRequestFields()
	l.WithFields(fields).Info("Request Successful")
}

func getCurrFuncName() string {
	return GetFuncName(4)
}

func getPrevFuncName() string {
	return GetFuncName(5)
}

func GetFuncName(numFrames int) string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(numFrames, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}
