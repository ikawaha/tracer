[![Build Status](https://travis-ci.org/ikawaha/tracer.svg?branch=master)](https://travis-ci.org/ikawaha/tracer)

HTTP request/response tracer middleware
---

The `Trace` middleware that keeps request copy and response records in context.
For example, this middleware can be used to record API requests and responses in the database. 
It is also useful when there is no HTTP request in the controller and 
only context can be obtained due to limitations of frameworks such as [Goa](https://goa.design).

If you do not need to record responses, you can use `RecordRequest` middleware.

|Middleware| Targets| Context Key| Context Object | Option |
|:---|:---|:---|:---|:---|
|Trace| Request and Response| TrackerKey | Tracker | DiscardResponseBody |
|RecordRequest| Request| RequestRecorderKey| RequestRecorder||


### Example

The `Trace` middleware used with [Goa](https://goa.design). Following is an example of a controller:

```go
func (s *calcsrvc) Add(ctx context.Context, p *calc.AddPayload) (res int, err error) {
	var header http.Header
	tracker, ok := ctx.Value(tracer.TrackerKey).(*tracer.Tracker)
	if ok {
		header = tracker.Request.Header
	}
	s.logger.Printf("calc.add, Header: %+v", header)
	return p.A + p.B, nil
}
```

**Blog**: http://ikawaha.hateblo.jp/entry/2019/12/05/235917

___

License MIT
