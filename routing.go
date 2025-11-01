package flowhttp

import (
	"fmt"
	"regexp"
	"strings"
)

// internal types representing streams and methods
type stream struct {
	steps []Step
	sink  Sink
}

type streamMethods struct {
	GET  *stream
	POST *stream
}

type dynamicStream struct {
	pattern       *regexp.Regexp
	methods       *streamMethods
	hasPathParams bool
}

// convertPathToRegex converts patterns like /user/:id or /files/*path to a named regex.
// Returns compiled regex and whether the pattern contains named params.
func convertPathToRegex(path string) (*regexp.Regexp, bool) {
	hasParams := false
	re := regexp.MustCompile(`:([a-zA-Z0-9_]+)`)
	if re.MatchString(path) {
		hasParams = true
	}
	replaced := re.ReplaceAllString(path, `(?P<$1>[^/]+)`)
	replaced = strings.ReplaceAll(replaced, "*", ".*")
	return regexp.MustCompile("^" + replaced + "$"), hasParams
}

// getStreamMethodsForPath resolves a path to either static or dynamic route.
// Returns streamMethods, extracted params (if any), or error when not found.
func (f *Flow) getStreamMethodsForPath(path string) (*streamMethods, map[string]string, error) {
	// static fast path
	if methods, exists := f.streams[path]; exists {
		return methods, nil, nil
	}
	// dynamic fallback (order preserved as registered)
	for _, d := range f.dynamicStreams {
		if d.pattern.MatchString(path) {
			params := make(map[string]string)
			if d.hasPathParams {
				matches := d.pattern.FindStringSubmatch(path)
				for i, name := range d.pattern.SubexpNames() {
					if i != 0 && name != "" {
						params[name] = matches[i]
					}
				}
			}
			return d.methods, params, nil
		}
	}
	return nil, nil, fmt.Errorf("no route found for path: %s", path)
}
