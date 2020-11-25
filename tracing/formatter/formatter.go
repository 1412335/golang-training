package main

import (
	"fmt"
	"golang-training/tracing/utils"
	"log"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

func main() {
	tracer, closer := utils.InitJaeger("formatter")
	defer closer.Close()

	http.HandleFunc("/formatv1", func(w http.ResponseWriter, r *http.Request) {
		// get span context from request
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		// span.kind=server
		span := tracer.StartSpan("format", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		// get baggage
		greeting := span.BaggageItem("greeting")
		if greeting == "" {
			greeting = "hello"
		}

		helloTo := r.FormValue("helloTo")
		helloStr := fmt.Sprintf("%s, %s!", greeting, helloTo)
		span.LogFields(
			otlog.String("event", "string-format"),
			otlog.String("value", helloStr),
		)
		w.Write([]byte(helloStr))
	})

	http.HandleFunc("/format", func(w http.ResponseWriter, r *http.Request) {
		// get span from request context
		span := opentracing.SpanFromContext(r.Context())
		// defer span.Finish() //=> duplicate span

		// get baggage
		greeting := span.BaggageItem("greeting")
		if greeting == "" {
			greeting = "hello"
		}

		helloTo := r.FormValue("helloTo")
		helloStr := fmt.Sprintf("%s, %s!", greeting, helloTo)
		// logs
		span.LogFields(
			otlog.String("event", "string-format"),
			otlog.String("value", helloStr),
		)
		// response
		w.Write([]byte(helloStr))
	})

	// log.Fatal(http.ListenAndServe(":8084", nil))
	log.Fatal(http.ListenAndServe("localhost:8084", nethttp.Middleware(tracer, http.DefaultServeMux)))
}
