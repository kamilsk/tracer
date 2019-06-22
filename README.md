> # 🧶 tracer
>
> Lightweight tracing mechanism.

## 💡 Idea

```go
func Do(ctx context.Context) {
	defer tracer.Fetch(ctx).Start().Stop()

	// do some job
}
```

Full description of the idea is available
[here](https://www.notion.so/octolab/tracer-098c6f9fe97b41dcac4a30074463dc8f?r=0b753cbf767346f5a6fd51194829a2f3).

## 🏆 Motivation

Coming soon.

## 🤼‍♂️ How to

```go
import "github.com/kamilsk/tracer"

func Do(ctx context.Context) {
	defer tracer.Fetch(ctx).Start().Mark("49cfe2b9-1942-47f1-92f6-6e7be7243845").Stop()

	// do some job

	tracer.Fetch(ctx).Breakpoint()

	// do some job

	tracer.Fetch(ctx).Breakpoint().Mark("c246ba1f-8a12-40ed-b4f7-b39289253ca1")
}
```

## 🧩 Integration

Coming soon.

---

made with ❤️ for everyone

[icon_build]:      https://travis-ci.org/kamilsk/tracer.svg?branch=master

[page_build]:      https://travis-ci.org/kamilsk/tracer
[page_promo]:      https://github.com/kamilsk/tracer
