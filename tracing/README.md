# Jaeger - Open tracing

## Testing
```
# jaeger
docker run \
  --rm \
  --name jaeger \
  -p6831:6831/udp \
  -p16686:16686 \
  jaegertracing/all-in-one:latest

# example
docker run \
  --rm \
  --link jaeger \
  --env JAEGER_AGENT_HOST=jaeger \
  --env JAEGER_AGENT_PORT=6831 \
  -p8080-8083:8080-8083 \
  jaegertracing/example-hotrod:latest \
  all

# frontend ui
http://127.0.0.1:8080/

# debug
http://127.0.0.1:8083/debug/vars

# jaeger ui
http://127.0.0.1:16686

# source
git clone https://github.com/jaegertracing/jaeger.git jaeger
cd jaeger/examples/hotrod
```

## Sampling Strategies
```
--sampling.strategies-file=/etc/jaeger/sampling_strategies.json
```

## Tutorial
- https://github.com/yurishkuro/opentracing-tutorial/tree/master/go

## Ref
- https://medium.com/opentracing/tracing-http-request-latency-in-go-with-opentracing-7cc1282a100a
- https://github.com/jaegertracing/jaeger/tree/master/examples/hotrod
- https://github.com/opentracing-contrib/go-stdlib
- https://github.com/opentracing/specification/blob/master/semantic_conventions.md

- https://www.jaegertracing.io/docs/1.17/sampling/


## Make span
```go
// from parent span
span := rootSpan.Tracer().StartSpan(
	"formatString",
	opentracing.ChildOf(rootSpan.Context()),
)
// from context
span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
defer span.Finish()

// get span
span := opentracing.SpanFromContext(req.Context())
```