package errors

import (
	"errors"
	"fmt"

	"github.com/rokmetro/logging-library/logutils"
)

//ErrorContext represents the context of an error message
type ErrorContext struct {
	message  string
	function string
}

//String converts the ErrorContext to a string
func (e ErrorContext) String() string {
	if e.function != "" {
		return fmt.Sprintf("%s() %s", e.function, e.message)
	}
	return e.message
}

//Error represents an error entity
type Error struct {
	root     *ErrorContext
	internal error
	tags     []string
	trace    []ErrorContext
}

//String returns the root message as a string
func (e *Error) String() string {
	return e.Root()
}

//Error returns the full trace context as an error
func (e *Error) Error() string {
	return e.TraceContext()
}

//Root returns the root message
func (e *Error) Root() string {
	if e == nil || e.root == nil {
		return ""
	}
	return e.root.message
}

//RootContext returns the root context
func (e *Error) RootContext() string {
	if e == nil || e.root == nil {
		return ""
	}
	root := e.root.String()
	if e.internal != nil {
		root += ":" + e.internal.Error()
	}
	return root
}

//Trace returns the trace messages
func (e *Error) Trace() string {
	if e == nil {
		return ""
	}
	trace := e.Root()
	for _, ctx := range e.trace {
		trace = traceError(ctx, trace)
	}
	return trace
}

//TraceContext returns the trace context
func (e *Error) TraceContext() string {
	if e == nil {
		return ""
	}
	trace := e.RootContext()
	for _, ctx := range e.trace {
		trace = traceContext(ctx, trace)
	}
	return trace
}

//RootErr returns the root message as an error
func (e *Error) RootErr() error {
	return errors.New(e.Root())
}

//RootContextErr returns the root context as an error
func (e *Error) RootContextErr() error {
	return errors.New(e.RootContext())
}

//TraceErr returns the trace messages as an error
func (e *Error) TraceErr() error {
	return errors.New(e.Trace())
}

//TraceContextErr returns the trace context as an error
func (e *Error) TraceContextErr() error {
	return errors.New(e.TraceContext())
}

//Tags returns the tags
func (e *Error) Tags() []string {
	if e == nil {
		return []string{}
	}
	return e.tags
}

//AddTag adds the provided tag and returns the result
func (e Error) AddTag(tag string) *Error {
	e.tags = append(e.tags, tag)
	return &e
}

//HasTag returns true if the Error has the provided tag
func (e *Error) HasTag(tag string) bool {
	if e == nil {
		return false
	}
	return logutils.ContainsString(e.tags, tag)
}

func (e Error) wrap(context *ErrorContext) *Error {
	if context == nil {
		return &e
	}
	if e.root == nil {
		e.root = context
	}
	e.trace = append(e.trace, *context)
	return &e
}

// func (e Error) setRoot(context *ErrorContext) *Error {
// 	if e.root != nil || context == nil {
// 		return &e
// 	}
// 	e.root = context
// 	return &e
// }

func traceError(ctx ErrorContext, trace string) string {
	return fmt.Sprintf("%s: %s", ctx.message, trace)
}

func traceContext(ctx ErrorContext, trace string) string {
	if ctx.function != "" {
		return fmt.Sprintf("%s() %s: [%s]", ctx.function, ctx.message, trace)
	}
	return fmt.Sprintf("%s: %s", ctx.message, trace)
}
