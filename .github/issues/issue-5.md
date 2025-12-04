---
id: 5
database_id: 491576484
node_id: MDU6SXNzdWU0OTE1NzY0ODQ=
status: open
title: "add possibility to preallocate checkpoints"
labels: ["type:feature","status:pending"]
url: https://github.com/kamilsk/tracer/issues/5
created_at: 2019-09-10T10:06:35Z
updated_at: 2019-09-17T07:51:46Z
---

# add possibility to preallocate checkpoints

```go
call := tracer.Fetch(ctx).With(make([]Checkpoint, 0, 10)).Start()
defer call.Stop()
...
call.Checkpoint()
...
call.Checkpoint("id", "labelX", "labelY")
...
```
