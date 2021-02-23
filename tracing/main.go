package main

import (
	"golang-training/tracing/pkg/cmd"
)

func main() {
	// cmd.Root(os.Args)
	cmd.Execute()
}

// type Service struct {
// 	tracer     opentracing.Tracer
// 	httpClient *tracing.HTTPClient
// }

// func (s *Service) formatString(ctx context.Context, helloTo string) string {
// 	// span := rootSpan.Tracer().StartSpan(
// 	// 	"formatString",
// 	// 	opentracing.ChildOf(rootSpan.Context()),
// 	// )
// 	// span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
// 	// defer span.Finish()

// 	// get child span
// 	// span := opentracing.SpanFromContext(req.Context())
// 	span := opentracing.SpanFromContext(ctx)

// 	// do request to formatter service
// 	v := url.Values{}
// 	v.Set("helloTo", helloTo)
// 	url := "http://localhost:8084/format?" + v.Encode()
// 	resp, err := s.httpClient.Do(ctx, url)
// 	if err != nil {
// 		// set span tag error=true
// 		ext.LogError(span, err)
// 		panic(err.Error())
// 	}

// 	helloStr := string(resp)

// 	// client logs span
// 	span.LogFields(
// 		log.String("event", "formatString"),
// 		log.String("value", helloStr),
// 	)

// 	return helloStr
// }

// func (s *Service) printHello(ctx context.Context, helloStr string) {
// 	// span := rootSpan.Tracer().StartSpan(
// 	// 	"printHello",
// 	// 	opentracing.ChildOf(rootSpan.Context()),
// 	// )

// 	// span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
// 	// defer span.Finish()

// 	// get root span
// 	// span := opentracing.SpanFromContext(req.Context())
// 	span := opentracing.SpanFromContext(ctx)

// 	v := url.Values{}
// 	v.Set("helloStr", helloStr)
// 	url := "http://localhost:8085/publish?" + v.Encode()
// 	_, err := s.httpClient.Do(ctx, url)
// 	if err != nil {
// 		// set span tag error=true
// 		ext.LogError(span, err)
// 		panic(err.Error())
// 	}

// 	span.LogKV("event", "printHello")
// }

// func main() {
// 	if len(os.Args) != 3 {
// 		panic("ERROR: Expecting 2 argument")
// 	}
// 	helloTo := os.Args[1]
// 	greeting := os.Args[2]

// 	logger, _ = zap.NewDevelopment(
// 		zap.AddStacktrace(zapcore.FatalLevel),
// 		zap.AddCallerSkip(1),
// 	)
// 	metricsFactory = jexpvar.NewFactory(10) // 10 buckets for histograms
// 	logger.Info("Using expvar as metrics backend")

// 	// init tracer
// 	// service name: hello-world
// 	tracer := tracing.Init("hello-world", metricsFactory, logs.NewFactory(logger))
// 	// defer closer.Close()
// 	// need to set with StartSpanFromContext
// 	opentracing.SetGlobalTracer(tracer)

// 	// start root-span
// 	// operation name: say-hello
// 	span := tracer.StartSpan("say-hello")
// 	// set tag
// 	span.SetTag("hello-to", helloTo)
// 	defer span.Finish()

// 	// baggage
// 	span.SetBaggageItem("greeting", greeting)

// 	// attach root-span to context & pass ctx to child services
// 	ctx := context.Background()
// 	ctx = opentracing.ContextWithSpan(ctx, span)

// 	// http
// 	client := &http.Client{Transport: &nethttp.Transport{}}
// 	httpClient := &tracing.HTTPClient{
// 		Client: client,
// 		Tracer: tracer,
// 	}

// 	// call child services
// 	service := &Service{
// 		tracer:     tracer,
// 		httpClient: httpClient,
// 	}
// 	helloStr := service.formatString(ctx, helloTo)
// 	service.printHello(ctx, helloStr)
// }
