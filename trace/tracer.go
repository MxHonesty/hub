package trace

// Interface for an object capable of tracing events.
type Tracer interface {
	Trace(...interface{})
}
