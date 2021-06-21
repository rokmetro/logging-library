# logging-library

Logging package for a standard logging interface across microservices.
Includes context information such as request ids, span ids, stack trace and other unstructured context data useful for tracking and debugging purposes.

```
type Log struct {
	*StandardLogger
	traceID    string
	spanID     string
	prevSpanID string
	stackTrace []string
	context    map[string]interface{}
}
```

- traceID- generated once for an incoming request.
- spanID- generated for each new API call.
- prevSpanID- spanID of the previous function
- stackTrace- list of function calls made so far
- context- map for other unstructured context data

If client sends a request to service 1, a trace-id and span-id is generated for service 1. If it makes any calls to other microservices, a span-id gets generated for each subsequent service, but trace-id stays the same.
