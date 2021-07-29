package trace

import (
	"fmt"
	"io"
)

// Interface for an object capable of tracing events.
type Tracer interface {
	Trace(...interface{})
}

// Returns a Tracer that outputs to an Writer interface.
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Returns a Tracer that will ignore calls to Trace.
func Off() Tracer {
	return &nilTracer{}
}
