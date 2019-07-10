> # 🧶 tracer
>
> Simple, lightweight tracing.

[![Build][icon_build]][page_build]
[![Quality][icon_quality]][page_quality]
[![Documentation][icon_docs]][page_docs]
[![Coverage][icon_coverage]][page_coverage]
[![Awesome][icon_awesome]][page_awesome]

## 💡 Idea

The tracer provides API to trace execution flow.

```go
func Do(ctx context.Context) {
	defer tracer.Fetch(ctx).Start().Stop()

	// do some heavy job
}
```

Full description of the idea is available [here][design].

## 🏆 Motivation

At [Avito](https://tech.avito.ru), we use the [Jaeger](https://www.jaegertracing.io) - a distributed tracing platform.
It is handy in most cases, but at production, we also use sampling. So, what is a problem, you say?

I had 0.02% requests with a `write: broken pipe` error and it was difficult to find the appropriate one in
the [Sentry](https://sentry.io) which also has trace related to it in the [Jaeger](https://www.jaegertracing.io).

For that reason, I wrote the simple solution to handle this specific case and found the bottleneck in our code quickly.

## 🤼‍♂️ How to

```go
import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/kamilsk/tracer"
)

func InjectTracer(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		req = req.WithContext(tracer.Inject(req.Context(), make([]*tracer.Call, 0, 10)))
		handler.ServeHTTP(rw, req)
	})
}

func Handle(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second)
	defer cancel()

	call := tracer.Fetch(req.Context()).Start(req.Header.Get("X-Request-Id"))
	defer call.Stop()

	...

	call.Checkpoint("serialize")
	data, err := FetchData(ctx, req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	call.Checkpoint("store")
	if err := StoreIntoDatabase(ctx, data); err != nil {
		http.Error(rw,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func FetchData(ctx context.Context, r io.Reader) (Data, error) {
	defer tracer.Fetch(ctx).Start().Stop()

	// fetch a data into a struct
}

func StoreIntoDatabase(ctx context.Context, data Data) error {
	defer tracer.Fetch(ctx).Start().Stop()

	// store the data into a database
}
```

Output:

```
allocates at call stack: 0, detailed call stack:
	call Handle [ca7a87c4-58d0-4fdf-857c-ef49fc3bf271]: 14.038083ms, allocates: 2
		checkpoint [serialize]: 1.163587ms
		checkpoint [store]: 2.436265ms
	call FetchData: 1.192829ms, allocates: 0
	call StoreIntoDatabase: 10.428663ms, allocates: 0
```

## 🧩 Integration

The library uses [SemVer](https://semver.org) for versioning, and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe through major releases.
You can use [go modules](https://github.com/golang/go/wiki/Modules) or
[dep](https://golang.github.io/dep/) to manage its version.

```bash
$ go get -u github.com/kamilsk/tracer

$ dep ensure -add github.com/kamilsk/tracer
```

---

made with ❤️ for everyone

[icon_awesome]:     https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg
[icon_build]:       https://travis-ci.org/kamilsk/tracer.svg?branch=master
[icon_coverage]:    https://api.codeclimate.com/v1/badges/fb66449d1f5c64542377/test_coverage
[icon_docs]:        https://godoc.org/github.com/kamilsk/tracer?status.svg
[icon_quality]:     https://goreportcard.com/badge/github.com/kamilsk/tracer

[page_awesome]:     https://github.com/avelino/awesome-go#performance
[page_build]:       https://travis-ci.org/kamilsk/tracer
[page_coverage]:    https://codeclimate.com/github/kamilsk/tracer/test_coverage
[page_docs]:        https://godoc.org/github.com/kamilsk/tracer
[page_quality]:     https://goreportcard.com/report/github.com/kamilsk/tracer

[design]:           https://www.notion.so/octolab/tracer-098c6f9fe97b41dcac4a30074463dc8f?r=0b753cbf767346f5a6fd51194829a2f3
[egg]:              https://github.com/kamilsk/egg
[promo]:            https://github.com/kamilsk/tracer

[tmp.docs]:         https://nicedoc.io/kamilsk/tracer?theme=dark
[tmp.history]:      https://github.githistory.xyz/kamilsk/tracer/blob/master/README.md
