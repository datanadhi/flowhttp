package flowhttp

import (
	"encoding/json"
	"io"
	"net/http"
)

type ctxKey struct{}

// package-local unique key for request context values
var paramsKey ctxKey = struct{}{}

// FlowContext carries request/response and per-request locals/params.
type FlowContext struct {
	Request  *http.Request
	Response http.ResponseWriter
	local    map[string]any
	Params   map[string]string
}

// Set, Get, Delete are helpers to store small local values.
func (f *FlowContext) Set(key string, value any) { f.local[key] = value }
func (f *FlowContext) Get(key string) any        { return f.local[key] }
func (f *FlowContext) Delete(key string)         { delete(f.local, key) }

// Param returns a named path parameter (empty string if missing).
func (f *FlowContext) Param(name string) string { return f.Params[name] }

// Sink is the user handler type. ServeHTTP builds FlowContext from *http.Request.
type Sink func(*FlowContext)

// ServeHTTP makes Flow compatible with Go’s http package.
//
// Go’s http server (http.ListenAndServe) expects any router or handler
// to implement the http.Handler interface, which requires a ServeHTTP
// method. This method is called automatically for each incoming request.
//
// You don’t need to call ServeHTTP directly — it’s used internally
// so Flow can act as a standard HTTP handler.
func (h Sink) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params map[string]string
	if val := r.Context().Value(paramsKey); val != nil {
		if p, ok := val.(map[string]string); ok {
			params = p
		}
	}
	ctx := &FlowContext{
		Response: w,
		Request:  r,
		local:    make(map[string]any),
		Params:   params,
	}
	h(ctx)
}

// JSON serializes the given data to JSON and writes it to the response.
// It automatically sets the correct Content-Type header and handles encoding errors.
func (f *FlowContext) JSON(status int, data any) {
	f.Response.Header().Set("Content-Type", "application/json")
	f.Response.WriteHeader(status)

	if err := json.NewEncoder(f.Response).Encode(data); err != nil {
		http.Error(f.Response, `{"error": "failed to encode JSON"}`, http.StatusInternalServerError)
	}
}

// BindJSON reads and parses JSON from the request body into the given struct/map.
func (f *FlowContext) BindJSON(v any) error {
	body, err := io.ReadAll(f.Request.Body)
	if err != nil {
		http.Error(f.Response, "failed to read request body", http.StatusBadRequest)
		return err
	}
	defer f.Request.Body.Close()

	if err := json.Unmarshal(body, v); err != nil {
		http.Error(f.Response, "invalid JSON", http.StatusBadRequest)
		return err
	}

	return nil
}
