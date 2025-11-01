package flowhttp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// ServeHTTP implements http.Handler and dispatches requests through Steps -> Sink.
func (f *Flow) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	streamMethods, params, err := f.getStreamMethodsForPath(path)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	// attach params to request context so Sink.ServeHTTP can pick them up
	if params != nil {
		req = req.WithContext(context.WithValue(req.Context(), paramsKey, params))
	}

	var s *stream
	switch method {
	case http.MethodGet:
		s = streamMethods.GET
	case http.MethodPost:
		s = streamMethods.POST
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if s == nil {
		http.NotFound(w, req)
		return
	}

	// build middleware chain (wrap in reverse)
	sink := s.sink
	for i := len(s.steps) - 1; i >= 0; i-- {
		sink = s.steps[i](sink)
	}

	// call the top-level sink (it will build FlowContext)
	sink.ServeHTTP(w, req)
}

// Run starts the HTTP server and supports graceful shutdown.
// port can be int, string (":8080" or "8080"), or nil (defaults to :8080).
func (f *Flow) Run(port any) error {
	addr := ":8080"
	switch v := port.(type) {
	case nil:
	case int:
		addr = fmt.Sprintf(":%d", v)
	case string:
		if v != "" {
			if v[0] != ':' {
				addr = ":" + v
			} else {
				addr = v
			}
		}
	default:
		return fmt.Errorf("invalid port type")
	}

	srv := &http.Server{Addr: addr, Handler: f}
	errChan := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown error: %v", err)
		}
		return nil
	case err := <-errChan:
		return fmt.Errorf("server error: %v", err)
	}
}
