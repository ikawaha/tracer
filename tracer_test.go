package tracer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Controller(t *testing.T, ctx context.Context) {
	t.Helper()
	tracker, ok := ctx.Value(TrackerKey).(*Tracker)
	if !ok {
		t.Fatal("tracker not found")
	}
	// URL
	if got, expected := tracker.Request.URL.Path, "/hello"; got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
	// headers
	if got, expected := tracker.Request.Header.Get("User-Agent"), "Gopher-Client"; got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
	if got, expected := tracker.Request.Header.Get("Content-Type"), "application/json"; got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func ResponseCheckMiddleware(t *testing.T) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			ctx := r.Context()
			tracker, ok := ctx.Value(TrackerKey).(*Tracker)
			if !ok {
				t.Fatal("tracker not found")
			}
			if got, expected := tracker.Response.Body.String(), "goodbye"; got != expected {
				t.Errorf("expected %q, got %q", expected, got)
			}
			if got, expected := tracker.Response.Status, http.StatusOK; got != expected {
				t.Errorf("expected %q, got %q", expected, got)
			}
		})
	}
}

func TestTracerMiddleware(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/hello",
		func(w http.ResponseWriter, r *http.Request) {
			Controller(t, r.Context())   // The controller knows context only.
			w.WriteHeader(http.StatusOK) // responses
			w.Write([]byte("goodbye"))
		},
	)
	// middlewares
	var handler http.Handler = mux
	handler = ResponseCheckMiddleware(t)(handler) // This middleware can know the response.
	handler = Trace(handler)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	req, err := http.NewRequest("POST", ts.URL+"/hello", nil)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	req.Header.Add("User-Agent", "Gopher-Client")
	req.Header.Add("Content-Type", "application/json")

	cl := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := cl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected OK, got %v", resp.Status)
	}
}
