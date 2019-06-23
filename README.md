> # 🧶 tracer
>
> Simple, lightweight tracing mechanism.

[![Documentation][icon_docs]][page_docs]

## 💡 Idea

```go
func Do(ctx context.Context) {
	defer tracer.Fetch(ctx).Start().Stop()

	// do some heavy job
}
```

Full description of the idea is available
[here](https://www.notion.so/octolab/tracer-098c6f9fe97b41dcac4a30074463dc8f?r=0b753cbf767346f5a6fd51194829a2f3).

## 🏆 Motivation

In [Avito](https://tech.avito.ru), we use the [Jaeger](https://www.jaegertracing.io) - a distributed tracing platform.
It is handy in most cases, but at production, we also use sampling. So, what is a problem, you say?

I had 0.02% requests with a `write: broken pipe` error and it was difficult to find the appropriate one in
the [Sentry](https://sentry.io) which also has trace related to it in the [Jaeger](https://www.jaegertracing.io).

For that reason, I wrote the simple solution to handle this specific case and found the bottleneck in our code quickly.

## 🤼‍♂️ How to

```go
import (
	"context"
	"net/http"
	"time"

	"github.com/kamilsk/tracer"
)

func InjectTracer(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		req = req.WithContext(tracer.Inject(req.Context(), make([]tracer.Call, 0, 10)))
		handler.ServeHTTP(rw, req)
	})
}

func Handle(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second)
	defer cancel()

	call := tracer.Fetch(req.Context()).Start().Mark(req.Header.Get("X-Request-Id"))
	defer call.Stop()

    ...

	call.Checkpoint().Mark("serialize")
	data := FetchData(ctx, req.Body)

	call.Checkpoint().Mark("store")
	if err := StoreIntoDatabase(ctx, data); err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func FetchData(ctx context.Context, r io.Reader) Data {
	defer tracer.Fetch(ctx).Start().Stop()

	// fetch a data into a struct
}

func StoreIntoDatabase(ctx context.Context, data Data) error {
	defer tracer.Fetch(ctx).Start().Stop()

	// store the data into a database
}
```

## 🧩 Integration

Coming soon.

---

made with ❤️ for everyone

[icon_build]:      https://travis-ci.org/kamilsk/tracer.svg?branch=master
[icon_docs]:       https://godoc.org/github.com/kamilsk/tracer?status.svg

[page_build]:      https://travis-ci.org/kamilsk/tracer
[page_docs]:       https://godoc.org/github.com/kamilsk/tracer
[page_promo]:      https://github.com/kamilsk/tracer
