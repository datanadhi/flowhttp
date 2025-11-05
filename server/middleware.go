package server

// Step is a middleware: it receives next Sink and returns a Sink.
type Step func(Sink) Sink

// CreateStep is a convenient wrapper so users can write (next, ctx) style middleware.
func CreateStep(fn func(next Sink, ctx *FlowContext)) Step {
	return func(next Sink) Sink {
		return func(ctx *FlowContext) { fn(next, ctx) }
	}
}
