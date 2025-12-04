---
id: 13
database_id: 777043394
node_id: MDU6SXNzdWU3NzcwNDMzOTQ=
status: open
title: "bucket to store debug info"
labels: ["type:feature","status:pending"]
url: https://github.com/kamilsk/tracer/issues/13
created_at: 2020-12-31T14:20:04Z
updated_at: 2020-12-31T14:20:11Z
---

# bucket to store debug info

```go
call.Store(Request{
  Type: HTTP
  ...
})

call.Store(Request{
  Type: SQL
  ...
})
```
