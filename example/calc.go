package calcapi

import (
	"context"
	"log"

	"calc/gen/calc"
	"github.com/ikawaha/tracer"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcsrvc struct {
	logger *log.Logger
}

// NewCalc returns the calc service implementation.
func NewCalc(logger *log.Logger) calc.Service {
	return &calcsrvc{logger}
}

// Add implements add.
func (s *calcsrvc) Add(ctx context.Context, p *calc.AddPayload) (res int, err error) {
	var header string
	tracker, ok := ctx.Value(tracer.TrackerKey).(*tracer.Tracker)
	if ok {
		header = tracker.Request.Header["Content-Type"][0]
	}
	s.logger.Printf("calc.add, Content-Type: %+v", header)
	return p.A + p.B, nil
}
