---
id: 4
database_id: 481041142
node_id: MDU6SXNzdWU0ODEwNDExNDI=
status: open
title: "example with sla"
labels: ["type:docs","status:pending"]
url: https://github.com/kamilsk/tracer/issues/4
created_at: 2019-08-15T08:06:15Z
updated_at: 2019-09-17T07:51:37Z
---

# example with sla

```go
func Middleware(handler http.Handler) http.Handler {
  return http.HandlerFunc(func (...) {
    ctx := context.WithValue(req.Context(), key, tracer)
    http.ServeHTTP(rw, req.WithContex(ctx))
    if ctx.Err() != nil {
      // fetch tracer, write trace into Sentry
    }
  })
}
```
