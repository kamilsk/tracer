package tracer_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kamilsk/tracer"
)

type Data struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

func Example() {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/message",
		strings.NewReader(`{"title": "tracer", "subtitle": "🧶 Simple, lightweight tracing mechanism."}`))
	req.Header.Set("X-Request-Id", "ca7a87c4-58d0-4fdf-857c-ef49fc3bf271")

	handler := InjectTracer(FlushTracer(http.HandlerFunc(Handle)))
	handler.ServeHTTP(rec, req)

	_, _ = io.Copy(os.Stdout, strings.NewReader(stabilize(rec.Body.String())))
	// Output:
	// allocates at call stack: 0, detailed call stack:
	// 	call Handle [ca7a87c4-58d0-4fdf-857c-ef49fc3bf271]: 12.345678ms, allocates: 2
	// 		checkpoint [serialize]: 1.234567ms
	// 		checkpoint [store]: 1.234567ms
	// 	call FetchData: 1.234567ms, allocates: 0
	// 	call StoreIntoDatabase: 12.345678ms, allocates: 0
}

func FlushTracer(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(rw, req)
		_, _ = rw.Write([]byte(tracer.Fetch(req.Context()).String()))
	})
}

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

	time.Sleep(time.Millisecond)

	call.Checkpoint("serialize")
	data, err := FetchData(ctx, req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	time.Sleep(time.Millisecond)

	call.Checkpoint("store")
	if err := StoreIntoDatabase(ctx, data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func FetchData(ctx context.Context, r io.Reader) (Data, error) {
	defer tracer.Fetch(ctx).Start().Stop()

	time.Sleep(time.Millisecond)
	var data Data
	err := json.NewDecoder(r).Decode(&data)
	return data, err
}

func StoreIntoDatabase(ctx context.Context, data Data) error {
	defer tracer.Fetch(ctx).Start().Stop()

	time.Sleep(10 * time.Millisecond)
	return gob.NewEncoder(bytes.NewBuffer(nil)).Encode(data)
}

func stabilize(raw string) string {
	raw = strings.Replace(raw, "tracer_test.", "", -1)
	raw = regexp.MustCompile(`Handle (.+): (\d{2}\.\d+ms)`).ReplaceAllString(raw, "Handle $1: 12.345678ms")
	raw = regexp.MustCompile(`\[serialize]: (\d\.\d+ms)`).ReplaceAllString(raw, "[serialize]: 1.234567ms")
	raw = regexp.MustCompile(`\[store]: (\d\.\d+ms)`).ReplaceAllString(raw, "[store]: 1.234567ms")
	raw = regexp.MustCompile(`FetchData: (\d\.\d+ms)`).ReplaceAllString(raw, "FetchData: 1.234567ms")
	raw = regexp.MustCompile(`StoreIntoDatabase: (\d{2}\.\d+ms)`).ReplaceAllString(raw, "StoreIntoDatabase: 12.345678ms")
	return raw
}
