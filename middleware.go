package tracer

import (
	"context"
	"net/http"
)

type constKey int

const (
	// TrackerKey is the request context key used to store the tracking data.
	TrackerKey constKey = iota + 1
	// RecordRequestKey is the request context key used to store the request data.
	RequestRecorderKey
)

type Option func(t *Tracker)

func DiscardResponseBody() Option {
	return func(t *Tracker) {
		t.Response.Discard = true
	}
}

// Trace is the middleware that traces HTTP requests and responses.
func Trace(options ...Option) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tracker := NewTracker(w, r)
			for _, v := range options {
				v(tracker)
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, TrackerKey, tracker)
			h.ServeHTTP(tracker, r.WithContext(ctx))
		})
	}
}

// RecordRequest is the middleware that store HTTP requests.
func RecordRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := NewRequestRecorder(r)
		ctx := r.Context()
		ctx = context.WithValue(ctx, RequestRecorderKey, &req)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
