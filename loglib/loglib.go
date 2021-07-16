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
	layer   int
}

//NewLog is a constructor for a log object
func (l *Logger) NewLog(traceID string, request RequestContext) *Log {
	if traceID == "" {
		traceID = uuid.New().String()
	}
	spanID := uuid.New().String()
	log := &Log{l, traceID, spanID, request, Fields{}, 0}
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

	log := &Log{l, traceID, spanID, request, Fields{}, 0}
	return log
}

func (l *Log) resetLayer() {
	l.layer = 0
}

func (l *Log) addLayer(layer int) {
	l.layer += layer
}

//getRequestFields() populates a map with all the fields of a request
//	layer: Number of function calls between caller and getRequestFields()
func (l *Log) getRequestFields() Fields {
	if l == nil {
		return Fields{}
	}

	fields := Fields{"trace_id": l.traceID, "span_id": l.spanID, "function_name": getLogPrevFuncName(l.layer)}
	l.resetLayer()

	return fields
}

//SetHeaders sets the trace and span id headers for a request to another service
//	This function should always be called when making a request to another rokwire service
func (l *Log) SetHeaders(r *http.Request) {
	if l == nil {
		return
	}

	r.Header.Set("trace-id", l.traceID)
	r.Header.Set("span-id", l.spanID)
}

//LogData logs and returns a data message at the designated level
//	level: The log level (Info, Debug, Warn, Error)
//	status: The status of the data
//	dataType: The data type
//	args: Any args that should be included in the message (nil if none)
func (l *Log) LogData(level logLevel, status logDataStatus, dataType logData, args logArgs) string {
	msg := DataMessage(status, dataType, args)
	l.addLayer(1)

	switch level {
	case Info:
		l.Error(msg)
	case Debug:
		l.Debug(msg)
	case Warn:
		l.Warn(msg)
	case Error:
		l.Error(msg)
	default:
		l.resetLayer()
	}

	return msg
}

//ErrorData logs and returns a data message for the given error
//	status: The status of the data
//	dataType: The data type
//	err: Error message
func (l *Log) ErrorData(status logDataStatus, dataType logData, err error) string {
	message := DataMessage(status, dataType, nil)

	l.addLayer(1)
	defer l.resetLayer()

	return l.LogError(message, err)
}

//RequestErrorData logs a data message and error and sets it as the HTTP response
//	w: The http response writer for the active request
//	status: The status of the data
//	dataType: The data type
//	err: The error received from the application
//	code: The HTTP response code to be set
//	hideDetails: Only provide 'msg' not 'err' in HTTP response when true
func (l *Log) RequestErrorData(w http.ResponseWriter, status logDataStatus, dataType logData, err error, code int, hideDetails bool) {
	message := DataMessage(status, dataType, nil)

	l.addLayer(1)
	defer l.resetLayer()

	l.RequestError(w, message, err, code, hideDetails)
}

//LogAction logs and returns an action message at the designated level
//	level: The log level (Info, Debug, Warn, Error)
//	status: The status of the action
//	action: The action that is occurring
//	dataType: The data type that the action is occurring on
//	args: Any args that should be included in the message (nil if none)
func (l *Log) LogAction(level logLevel, status logActionStatus, action logAction, dataType logData, args logArgs) string {
	msg := ActionMessage(status, action, dataType, args)
	l.addLayer(1)

	switch level {
	case Info:
		l.Error(msg)
	case Debug:
		l.Debug(msg)
	case Warn:
		l.Warn(msg)
	case Error:
		l.Error(msg)
	default:
		l.resetLayer()
	}

	return msg
}

//ErrorAction logs and returns an action message for the given error
//	action: The action that is occurring
//	dataType: The data type that the action is occurring on
//	err: Error message
func (l *Log) ErrorAction(action logAction, dataType logData, err error) string {
	message := ActionMessage(ErrorStatus, action, dataType, nil)

	l.addLayer(1)
	defer l.resetLayer()

	return l.LogError(message, err)
}

//RequestErrorAction logs an action message and error and sets it as the HTTP response
//	w: The http response writer for the active request
//	action: The action that is occurring
//	dataType: The data type
//	err: The error received from the application
//	code: The HTTP response code to be set
//	hideDetails: Only provide 'msg' not 'err' in HTTP response when true
func (l *Log) RequestErrorAction(w http.ResponseWriter, action logAction, dataType logData, err error, code int, hideDetails bool) {
	message := ActionMessage(ErrorStatus, action, dataType, nil)

	l.addLayer(1)
	defer l.resetLayer()

	l.RequestError(w, message, err, code, hideDetails)
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
	msg := fmt.Sprintf("%s: %v", message, err)
	if l == nil || l.logger == nil {
		return msg
	}

	requestFields := l.getRequestFields()
	requestFields["error"] = err
	l.logger.withFields(requestFields).Error(message)
	return msg
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

//RequestError logs the provided message and error and sets it as the HTTP response
//	Params:
//		w: The http response writer for the active request
//		msg: The error message
//		err: The error received from the application
//		code: The HTTP response code to be set
//		hideDetails: Only provide 'msg' not 'err' in HTTP response when true
func (l *Log) RequestError(w http.ResponseWriter, msg string, err error, code int, hideDetails bool) {
	l.addLayer(1)
	defer l.resetLayer()

	l.SetContext("status_code", code)

	msg = fmt.Sprintf("%d - %s", code, msg)
	detailMsg := l.LogError(msg, err)
	if !hideDetails {
		msg = detailMsg
	}
	http.Error(w, msg, code)
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
