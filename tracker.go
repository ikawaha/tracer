package tracer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// RequestRecorder reads http request and records it's body.
type RequestRecorder struct {
	*http.Request
	Payload json.RawMessage
}

// RequestRecorder returns copied request.
func NewRequestRecorder(r *http.Request) RequestRecorder {
	b, err := ioutil.ReadAll(r.Body)
	if err == nil {
		_ = r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewReader(b))
	}
	return RequestRecorder{
		Request: r,
		Payload: b,
	}
}

// ResponseRecorder records response status, body and error code.
type ResponseRecorder struct {
	// ErrorCode is the code of the error returned by the action if any.
	ErrorCode string
	// Status is the response HTTP status code.
	Status int
	// Body is the response body
	Body bytes.Buffer
}

// Tracker implements http.ResponsWriter
var _ http.ResponseWriter = (*Tracker)(nil)

// Tracker represents data tracking requests and responses.
type Tracker struct {
	// Request
	Request RequestRecorder
	// Response
	Response ResponseRecorder
	http.ResponseWriter
}

// NewTracker creates http.ResponseWrite that records request and response.
func NewTracker(w http.ResponseWriter, r *http.Request) *Tracker {
	return &Tracker{
		Request:        NewRequestRecorder(r),
		ResponseWriter: w,
	}
}

// Written returns true if the response was written, false otherwise.
func (r *Tracker) Written() bool {
	return r.Response.Status != 0
}

// WriteHeader records the response status code and calls the underlying writer.
func (r *Tracker) WriteHeader(status int) {
	r.Response.Status = status
	r.ResponseWriter.WriteHeader(status)
}

// Write records the amount of data written and calls the underlying writer.
func (r *Tracker) Write(b []byte) (int, error) {
	if !r.Written() {
		r.WriteHeader(http.StatusOK)
	}
	r.Response.Body.Write(b)
	return r.ResponseWriter.Write(b)
}
