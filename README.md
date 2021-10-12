# WARNING - MIGRATED TO [ROKWIRE](https://github.com/rokwire/logging-library-go)
This library has been moved to the ROKWIRE open-source here: https://github.com/rokwire/logging-library-go. Please update the dependencies in all projects using this library to reflect this change.

All future work will be done in the ROKWIRE organization, so please submit issues and pull requests on the new repository.

# logging-library

Logging package for a standard logging interface across microservices.
Includes context information such as request ids, span ids, stack trace and other unstructured context data useful for tracking and debugging purposes.

```
type Log struct {
	logger  *Logger
	traceID string
	spanID  string
	request RequestContext
	context Fields
}
```

- traceID - generated once for an incoming request
- spanID  - generated for each new API call
- request - details about the request
- context - map for other unstructured context data

If client sends a request to service 1, a trace-id and span-id is generated for service 1. If it makes any calls to other microservices, a span-id gets generated for each subsequent service, but trace-id stays the same.

To get started, take a look at `example/app.go`
