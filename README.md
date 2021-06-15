# logwrapper

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

- traceID- gets generated at the top level for each incoming request.
- spanID- gets generated for each function call while servicing a request.
- prevSpanID- spanID of the previous function
- stackTrace- list of function calls made so far
- context- map for other unstructured context data
