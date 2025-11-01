package flowhttp

import (
	"fmt"
	"strings"
)

type Branch struct {
	path  string
	steps []Step
	flow  *Flow
}

// Flow is the top-level router object.
type Flow struct {
	streams        map[string]*streamMethods
	dynamicStreams []dynamicStream
	Branch
}

// NewFlow creates a root Flow.
func NewFlow() *Flow {
	f := &Flow{streams: make(map[string]*streamMethods)}
	f.flow = f
	return f
}

// Fork creates a sub-branch with a path prefix and inherited steps.
func (b *Branch) Fork(path string, steps []Step) *Branch {
	if path == "/" {
		path = ""
	}
	return &Branch{
		path:  b.path + path,
		steps: append(b.steps, steps...),
		flow:  b.flow,
	}
}

// ClearSteps clears inherited steps for this branch.
func (b *Branch) ClearSteps() *Branch {
	b.steps = nil
	return b
}

// Stream registers a route handler for method+path under this branch.
func (b *Branch) Stream(method string, path string, steps []Step, sink Sink) {
	finalPath := b.path + path
	finalSteps := append(b.steps, steps...)

	f := b.flow
	if f.streams == nil {
		f.streams = make(map[string]*streamMethods)
	}
	m := f.streams[finalPath]
	if m == nil {
		m = &streamMethods{}
	}
	h := &stream{steps: finalSteps, sink: sink}

	switch method {
	case "GET":
		m.GET = h
	case "POST":
		m.POST = h
	default:
		panic(fmt.Errorf("unsupported http method %s", method))
	}

	// dynamic route detection uses original path fragment (not prefixed finalPath),
	// so we check 'path' for params/wildcards to keep intent clear.
	if strings.Contains(path, ":") || strings.Contains(path, "*") {
		pattern, hasParams := convertPathToRegex(finalPath) // store compiled regex using finalPath
		f.dynamicStreams = append(f.dynamicStreams, dynamicStream{pattern, m, hasParams})
	} else {
		f.streams[finalPath] = m
	}
}
