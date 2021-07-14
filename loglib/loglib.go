package loglib

import (
	"net/http"
	"runtime"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Fields map[string]interface{}

func (f Fields) ToMap() map[string]interface{} {
	return f
}

//StandardLogger struct defines a wrapper for a logger object
type StandardLogger struct {
	entry *logrus.Entry
}

func (l *StandardLogger) withFields(fields Fields) *StandardLogger {
	return &StandardLogger{entry: l.entry.WithFields(fields.ToMap())}
}

//Fatal prints the log with a fatal error message and stops the service instance
//WARNING: Please only use for critical error messages that should prevent the service from running
func (l *StandardLogger) Fatal(message string) {
	l.entry.Fatal(message)
}

//Fatalf prints the log with a fatal format error message and stops the service instance
//WARNING: Please only use for critical error messages that should prevent the service from running
func (l *StandardLogger) Fatalf(message string, args ...interface{}) {
	l.entry.Fatalf(message, args)
}

//Error prints the log at error level with given message
func (l *StandardLogger) Error(message string) {
	l.entry.Error(message)
}

//ErrorWithFields prints the log at error level with given fields and message
func (l *StandardLogger) ErrorWithFields(message string, fields Fields) {
	l.entry.WithFields(fields.ToMap()).Error(message)
}

//Errorf prints the log at error level with given formatted string
func (l *StandardLogger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args)
}

//Info prints the log at info level with given message
func (l *StandardLogger) Info(message string) {
	l.entry.Info(message)
}

//InfoWithFields prints the log at info level with given fields and message
func (l *StandardLogger) InfoWithFields(message string, fields Fields) {
	l.entry.WithFields(fields.ToMap()).Info(message)
}

//Infof prints the log at info level with given formatted string
func (l *StandardLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args)
}

//Debug prints the log at debug level with given message
func (l *StandardLogger) Debug(message string) {
	l.entry.Debug(message)
}

//DebugWithFields prints the log at debug level with given fields and message
func (l *StandardLogger) DebugWithFields(message string, fields Fields) {
	l.entry.WithFields(fields.ToMap()).Debug(message)
}

//Debugf prints the log at debug level with given formatted string
func (l *StandardLogger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args)
}

//Warn prints the log at warn level with given message
func (l *StandardLogger) Warn(message string) {
	l.entry.Warn(message)
}

//WarnWithFields prints the log at warn level with given fields and message
func (l *StandardLogger) WarnWithFields(message string, fields Fields) {
	l.entry.WithFields(fields.ToMap()).Warn(message)
}

//Warnf prints the log at warn level with given formatted string
func (l *StandardLogger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args)
}

//Log struct defines a log object of a request
type Log struct {
	logger     *StandardLogger
	traceID    string
	spanID     string
	prevSpanID string
	context    Fields
}

//NewLogger is constructor for a logger object with initial configuration at the service level
func NewLogger(serviceName string) *StandardLogger {
	var baseLogger = logrus.New()
	baseLogger.Formatter = &logrus.JSONFormatter{}
	standardFields := logrus.Fields{"serviceName": serviceName} //All common fields for logs of a given service
	contextLogger := &StandardLogger{baseLogger.WithFields(standardFields)}
	return contextLogger
}

//NewLog is a constructor for a log object for a request
func (l *StandardLogger) NewLog(traceID string, prevSpanID string) *Log {
	if traceID == "" {
		traceID = uuid.New().String()
	}
	spanID := uuid.New().String()
	log := &Log{l, traceID, spanID, prevSpanID, Fields{}}
	return log
}

//NewRequestLog is a constructor for a log object for a request
func (l *StandardLogger) NewRequestLog(r *http.Request) *Log {
	traceID := r.Header.Get("trace-id")
	if traceID == "" {
		traceID = uuid.New().String()
	}
	prevSpanID := r.Header.Get("span-id")
	spanID := uuid.New().String()
	log := &Log{l, traceID, spanID, prevSpanID, Fields{}}
	return log
}

//getRequestFields() populates a map with all the fields of a request
func (l *Log) getRequestFields() Fields {
	fields := Fields{"trace_id": l.traceID, "span_id": l.spanID, "function_name": getPrevFuncName()}
	return fields
}

func (l *Log) SetHeaders(r *http.Request) {
	r.Header.Set("trace-id", l.traceID)
	r.Header.Set("span-id", l.spanID)
}

//InvalidArg is a standard error interface for invalid arguments
func (l *Log) InvalidArg(argumentName string, argumentValue interface{}) {
	requestFields := l.getRequestFields()
	requestFields["argument"] = argumentName
	requestFields["value"] = argumentValue
	l.logger.withFields(requestFields).Error("Invalid argument")
}

// MissingArg is a standard error interface for missing arguments
func (l *Log) MissingArg(argumentName string) {
	requestFields := l.getRequestFields()
	requestFields["argument"] = argumentName
	l.logger.withFields(requestFields).Error("Missing argument")
}

//Info prints the log at info level with given message
func (l *Log) Info(message string) {
	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Info(message)
}

//InfoWithDetails prints the log at info level with given fields and message
func (l *Log) InfoWithDetails(message string, details Fields) {
	requestFields := l.getRequestFields()
	requestFields["details"] = details
	l.logger.withFields(requestFields).Info(message)
}

//Infof prints the log at info level with given formatted string
func (l *Log) Infof(format string, args ...interface{}) {
	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Infof(format, args)
}

//Debug prints the log at debug level with given message
func (l *Log) Debug(message string) {
	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Debug(message)
}

//DebugWithDetails prints the log at debug level with given fields and message
func (l *Log) DebugWithDetails(message string, details Fields) {
	requestFields := l.getRequestFields()
	requestFields["details"] = details
	l.logger.withFields(requestFields).Debug(message)
}

//Debugf prints the log at debug level with given formatted string
func (l *Log) Debugf(format string, args ...interface{}) {
	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Debugf(format, args)
}

//Warn prints the log at warn level with given message
func (l *Log) Warn(message string) {
	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Warn(message)
}

//WarnWithDetails prints the log at warn level with given details and message
func (l *Log) WarnWithDetails(message string, details Fields) {
	requestFields := l.getRequestFields()
	requestFields["details"] = details
	l.logger.withFields(requestFields).Warn(message)
}

//Warnf prints the log at warn level with given formatted string
func (l *Log) Warnf(format string, args ...interface{}) {
	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Warnf(format, args)
}

//Error prints the log at error level with given message
func (l *Log) Error(format string) {
	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Error(format)
}

//ErrorWithDetails prints the log at error level with given details and message
func (l *Log) ErrorWithDetails(message string, details Fields) {
	requestFields := l.getRequestFields()
	requestFields["details"] = details
	l.logger.withFields(requestFields).Error(message)
}

//Errorf prints the log at error level with given formatted string
func (l *Log) Errorf(format string, args ...interface{}) {
	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Errorf(format, args)
}

//TODO: More error interfaces to be added

//AddContext adds any relevant unstructured data to context map
func (l *Log) AddContext(fieldName string, value interface{}) {
	l.context[fieldName] = value
}

//PrintContext prints the entire context of a log object
func (l *Log) PrintContext() {
	fields := l.getRequestFields()
	fields["prev_span_id"] = l.prevSpanID
	fields["context"] = l.context
	l.logger.InfoWithFields("Request Complete", fields)
}

//getCurrFuncName- fetches the current function name
func getCurrFuncName() string {
	return GetFuncName(4)
}

//getPrevFuncName- fetches the previous function name
func getPrevFuncName() string {
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
