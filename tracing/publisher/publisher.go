package main

import (
	"golang-training/tracing/utils"
	"log"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
)

func main() {
	tracer, closer := utils.InitJaeger("publisher")
	defer closer.Close()

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		// get span from request context
		span := opentracing.SpanFromContext(r.Context())
		// defer span.Finish() //=> duplicate span

		helloStr := r.FormValue("helloStr")
		println(helloStr)

		// logs
		span.LogKV("event", "println")
	})

	// log.Fatal(http.ListenAndServe(":8085", nil))
	log.Fatal(http.ListenAndServe("localhost:8085", nethttp.Middleware(tracer, http.DefaultServeMux)))
}
