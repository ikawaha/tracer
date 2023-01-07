package calcapi

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"calc/gen/calc"
	"calc/gen/http/calc/server"
	"github.com/ikawaha/goahttpcheck"
	"github.com/ikawaha/tracer"
)

func TestCalcsrvc_Add(t *testing.T) {
	var b bytes.Buffer
	logger := log.New(&b, "", log.LstdFlags)

	checker := goahttpcheck.New()
	checker.Use(tracer.Trace())

	checker.Mount(
		server.NewAddHandler,
		server.MountAddHandler,
		calc.NewEndpoints(NewCalc(logger)).Add,
	)

	checker.Test(t, "GET", "/add/1/2").
		WithHeader("Content-Type", "application/json").
		Check()

	// log check
	if got, expected := b.String(), "calc.add, Content-Type: application/json\n"; !strings.HasSuffix(got, expected) {
		t.Errorf("expected suffix, %q, got %q", expected, got)
	}
}
