package loglib

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Fields map[string]interface{}

func (f Fields) ToMap() map[string]interface{} {
	return f
}

//Logger struct defines a wrapper for a logger object
type Logger struct {
	entry            *logrus.Entry
	sensitiveHeaders []string
}

//LoggerOpts provides configuration options for the Logger type
type LoggerOpts struct {
	//JsonFmt When true, logs will be output in JSON format. Otherwise logs will be in logfmt
	JsonFmt bool
	//SensitiveHeaders: A list of any headers that contain sensitive information and should not be logged
	//				    Defaults: Authorization, Csrf
	SensitiveHeaders []string
}

//NewLogger is constructor for a logger object with initial configuration at the service level
// Params:
//		serviceName: A meaningful service name to be associated with all logs
//		opts: Configuration options for the Logger
func NewLogger(serviceName string, opts *LoggerOpts) *Logger {
	var baseLogger = logrus.New()
	sensitiveHeaders := []string{"Authorization", "Csrf"}

	if opts != nil {
		if opts.JsonFmt {
			baseLogger.Formatter = &logrus.JSONFormatter{}
		} else {
			baseLogger.Formatter = &logrus.TextFormatter{}
		}

		sensitiveHeaders = append(sensitiveHeaders, opts.SensitiveHeaders...)
	}

	standardFields := logrus.Fields{"service_name": serviceName} //All common fields for logs of a given service
	contextLogger := &Logger{entry: baseLogger.WithFields(standardFields), sensitiveHeaders: sensitiveHeaders}
	return contextLogger
}

func (l *Logger) withFields(fields Fields) *Logger {
	return &Logger{entry: l.entry.WithFields(fields.ToMap())}
}

//Fatal prints the log with a fatal error message and stops the service instance
//WARNING: Please only use for critical error messages that should prevent the service from running
func (l *Logger) Fatal(message string) {
	l.entry.Fatal(message)
}

//Fatalf prints the log with a fatal format error message and stops the service instance
//WARNING: Please only use for critical error messages that should prevent the service from running
func (l *Logger) Fatalf(message string, args ...interface{}) {
	l.entry.Fatalf(message, args...)
}

//Error prints the log at error level with given message
func (l *Logger) Error(message string) {
	l.entry.Error(message)
}

//ErrorWithFields prints the log at error level with given fields and message
func (l *Logger) ErrorWithFields(message string, fields Fields) {
	l.entry.WithFields(fields.ToMap()).Error(message)
}

//Errorf prints the log at error level with given formatted string
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

//Info prints the log at info level with given message
func (l *Logger) Info(message string) {
	l.entry.Info(message)
}

//InfoWithFields prints the log at info level with given fields and message
func (l *Logger) InfoWithFields(message string, fields Fields) {
	l.entry.WithFields(fields.ToMap()).Info(message)
}

//Infof prints the log at info level with given formatted string
func (l *Logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

//Debug prints the log at debug level with given message
func (l *Logger) Debug(message string) {
	l.entry.Debug(message)
}

//DebugWithFields prints the log at debug level with given fields and message
func (l *Logger) DebugWithFields(message string, fields Fields) {
	l.entry.WithFields(fields.ToMap()).Debug(message)
}

//Debugf prints the log at debug level with given formatted string
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

//Warn prints the log at warn level with given message
func (l *Logger) Warn(message string) {
	l.entry.Warn(message)
}

//WarnWithFields prints the log at warn level with given fields and message
func (l *Logger) WarnWithFields(message string, fields Fields) {
	l.entry.WithFields(fields.ToMap()).Warn(message)
}

//Warnf prints the log at warn level with given formatted string
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

type RequestContext struct {
	Method     string
	Path       string
	Headers    map[string][]string
	PrevSpanID string
}

func (r RequestContext) String() string {
	return fmt.Sprintf("%s %s prev_span_id: %s headers: %v", r.Method, r.Path, r.PrevSpanID, r.Headers)
}

//Log struct defines a log object of a request
type Log struct {
	logger  *Logger
	traceID string
	spanID  string
	request RequestContext
	context Fields
}

//NewLog is a constructor for a log object
func (l *Logger) NewLog(traceID string, request RequestContext) *Log {
	if traceID == "" {
		traceID = uuid.New().String()
	}
	spanID := uuid.New().String()
	log := &Log{l, traceID, spanID, request, Fields{}}
	return log
}

//NewRequestLog is a constructor for a log object for a request
func (l *Logger) NewRequestLog(r *http.Request) *Log {
	if r == nil {
		return &Log{logger: l}
	}

	traceID := r.Header.Get("trace-id")
	if traceID == "" {
		traceID = uuid.New().String()
	}

	prevSpanID := r.Header.Get("span-id")
	spanID := uuid.New().String()

	method := r.Method
	path := r.URL.Path

	headers := make(map[string][]string)
	for key, value := range r.Header {
		var logValue []string
		//do not log sensitive information
		if containsString(l.sensitiveHeaders, key) {
			logValue = append(logValue, "---")
		} else {
			logValue = value
		}
		headers[key] = logValue
	}

	request := RequestContext{Method: method, Path: path, Headers: headers, PrevSpanID: prevSpanID}

	log := &Log{l, traceID, spanID, request, Fields{}}
	return log
}

//getRequestFields() populates a map with all the fields of a request
func (l *Log) getRequestFields() Fields {
	if l == nil {
		return Fields{}
	}

	fields := Fields{"trace_id": l.traceID, "span_id": l.spanID, "function_name": getRequestFieldsPrevFuncName()}
	return fields
}

func (l *Log) SetHeaders(r *http.Request) {
	if l == nil {
		return
	}

	r.Header.Set("trace-id", l.traceID)
	r.Header.Set("span-id", l.spanID)
}

//InvalidArg is a standard error interface for invalid arguments
func (l *Log) InvalidArg(argumentName string, argumentValue interface{}) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	requestFields["argument"] = argumentName
	requestFields["value"] = argumentValue
	l.logger.withFields(requestFields).Error("Invalid argument")
}

// MissingArg is a standard error interface for missing arguments
func (l *Log) MissingArg(argumentName string) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	requestFields["argument"] = argumentName
	l.logger.withFields(requestFields).Error("Missing argument")
}

//Info prints the log at info level with given message
func (l *Log) Info(message string) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Info(message)
}

//InfoWithDetails prints the log at info level with given fields and message
func (l *Log) InfoWithDetails(message string, details Fields) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	requestFields["details"] = details
	l.logger.withFields(requestFields).Info(message)
}

//Infof prints the log at info level with given formatted string
func (l *Log) Infof(format string, args ...interface{}) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Infof(format, args...)
}

//Debug prints the log at debug level with given message
func (l *Log) Debug(message string) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Debug(message)
}

//DebugWithDetails prints the log at debug level with given fields and message
func (l *Log) DebugWithDetails(message string, details Fields) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	requestFields["details"] = details
	l.logger.withFields(requestFields).Debug(message)
}

//Debugf prints the log at debug level with given formatted string
func (l *Log) Debugf(format string, args ...interface{}) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Debugf(format, args...)
}

//Warn prints the log at warn level with given message
func (l *Log) Warn(message string) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Warn(message)
}

//WarnWithDetails prints the log at warn level with given details and message
func (l *Log) WarnWithDetails(message string, details Fields) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	requestFields["details"] = details
	l.logger.withFields(requestFields).Warn(message)
}

//Warnf prints the log at warn level with given formatted string
func (l *Log) Warnf(format string, args ...interface{}) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Warnf(format, args...)
}

//LogError prints the log at error level with given message and error
//	Returns error message as string
func (l *Log) LogError(message string, err error) string {
	if l == nil || l.logger == nil {
		return ""
	}

	requestFields := l.getRequestFields()
	requestFields["error"] = err
	l.logger.withFields(requestFields).Error(message)
	return fmt.Sprintf("%s: %v", message, err)
}

//Error prints the log at error level with given message
// Note: If possible, use LogError() instead
func (l *Log) Error(message string) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Error(message)
}

//ErrorWithDetails prints the log at error level with given details and message
func (l *Log) ErrorWithDetails(message string, details Fields) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	requestFields["details"] = details
	l.logger.withFields(requestFields).Error(message)
}

//Errorf prints the log at error level with given formatted string
// Note: If possible, use LogError() instead
func (l *Log) Errorf(format string, args ...interface{}) {
	if l == nil || l.logger == nil {
		return
	}

	requestFields := l.getRequestFields()
	l.logger.withFields(requestFields).Errorf(format, args...)
}

//TODO: More error interfaces to be added

//AddContext adds any relevant unstructured data to context map
// If the provided key already exists in the context, an error is returned
func (l *Log) AddContext(fieldName string, value interface{}) error {
	if l == nil {
		return fmt.Errorf("error adding context: nil log")
	}

	if _, ok := l.context[fieldName]; ok {
		return fmt.Errorf("error adding context: %s already exists", fieldName)
	}

	l.context[fieldName] = value
	return nil
}

//SetContext sets the provided context key to the provided value
func (l *Log) SetContext(fieldName string, value interface{}) {
	l.context[fieldName] = value
}

//RequestReceived prints the request context of a log object
func (l *Log) RequestReceived() {
	if l == nil || l.logger == nil {
		return
	}

	fields := l.getRequestFields()
	fields["request"] = l.request
	l.logger.InfoWithFields("Request Received", fields)
}

//RequestComplete prints the context of a log object
func (l *Log) RequestComplete() {
	if l == nil || l.logger == nil {
		return
	}

	fields := l.getRequestFields()
	fields["context"] = l.context
	l.logger.InfoWithFields("Request Complete", fields)
}
