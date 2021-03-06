## log

```golang
package main

import (
	"context"
	"net/http"

	"github.com/chanxuehong/log"
	"github.com/chanxuehong/log/trace"
)

// set defaults if necessary.
// normal at the moment of load config.
func init() {
	log.SetFormatter(log.TextFormatter)
	log.SetLevelString(log.DebugLevelString)
	log.SetDefaultOptions([]log.Option{
		log.WithFormatter(log.TextFormatter),
		log.WithLevelString(log.DebugLevelString),
	})
}

type mockResponseWriter struct {
	http.ResponseWriter
	header http.Header
}

func (w *mockResponseWriter) Header() http.Header {
	return w.header
}

func main() {
	w := &mockResponseWriter{
		header: make(http.Header),
	}
	req, _ := http.NewRequest(http.MethodGet, "http://www.example.com", nil)
	httpHandler(w, req)
}

func httpHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		log.Debug("debug", "trace_id", w.Header().Get(trace.TraceIdHeaderKey))
	}()

	// In general, it is a middleware
	{
		// gets traceId from request header, if not exists, generates it.
		traceId, ok := trace.FromRequest(req)
		if !ok {
			traceId = trace.NewTraceId()
			req.Header.Set(trace.TraceIdHeaderKey, traceId)
		}
		// sets traceId to response header
		w.Header().Set(trace.TraceIdHeaderKey, traceId)
		// adds traceId to request.Context
		req = req.WithContext(trace.NewContext(req.Context(), traceId))
		// adds log.Logger to request.Context
		req = req.WithContext(log.NewContext(req.Context(), log.New(log.WithTraceId(traceId))))
	}

	l := log.MustFromRequest(req)
	l.Info("1.info message")
	l.Info("2.info message", "key1", 1, "key2", 2)

	l = l.WithField("key3", 3)
	l = l.WithFields("key4", 4, "key5", 5)
	l.Info("3.info message")

	ctx := log.NewContext(req.Context(), l)
	fn1(ctx)
}

func fn1(ctx context.Context) {
	l := log.MustFromContext(ctx)
	l.Info("4.info message")
	l.Info("5.info message", "key6", 6)

	l = l.WithField("key7", 7)
	ctx = log.NewContext(ctx, l)
	fn2(ctx)
}

func fn2(ctx context.Context) {
	l := log.MustFromContext(ctx)
	l.Info("6.info message", "key8", 8)

	fn3(ctx)
}

func fn3(ctx context.Context) {
	l := log.MustFromContext(ctx)
	l.Info("7.info message", "key9", 9)
}
```

```Text
time=2018-05-25 23:59:16.066, level=info, request_id=934d159e603411e8a1cab4d5bdb21e16, location=main.httpHandler(test1/main.go:61), msg=1.info message
time=2018-05-25 23:59:16.066, level=info, request_id=934d159e603411e8a1cab4d5bdb21e16, location=main.httpHandler(test1/main.go:62), msg=2.info message, key1=1, key2=2
time=2018-05-25 23:59:16.066, level=info, request_id=934d159e603411e8a1cab4d5bdb21e16, location=main.httpHandler(test1/main.go:66), msg=3.info message, key3=3, key4=4, key5=5
time=2018-05-25 23:59:16.066, level=info, request_id=934d159e603411e8a1cab4d5bdb21e16, location=main.fn1(test1/main.go:74), msg=4.info message, key3=3, key4=4, key5=5
time=2018-05-25 23:59:16.066, level=info, request_id=934d159e603411e8a1cab4d5bdb21e16, location=main.fn1(test1/main.go:75), msg=5.info message, key3=3, key4=4, key5=5, key6=6
time=2018-05-25 23:59:16.067, level=info, request_id=934d159e603411e8a1cab4d5bdb21e16, location=main.fn2(test1/main.go:84), msg=6.info message, key3=3, key4=4, key5=5, key7=7, key8=8
time=2018-05-25 23:59:16.067, level=info, request_id=934d159e603411e8a1cab4d5bdb21e16, location=main.fn3(test1/main.go:91), msg=7.info message, key3=3, key4=4, key5=5, key7=7, key9=9
time=2018-05-25 23:59:16.067, level=debug, request_id=, location=main.httpHandler.func1(test1/main.go:41), msg=debug, trace_id=934d159e603411e8a1cab4d5bdb21e16
```
