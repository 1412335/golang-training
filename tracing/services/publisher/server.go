package publisher

import (
	"context"
	"golang-training/tracing/pkg/config"
	"golang-training/tracing/pkg/log"
	"golang-training/tracing/pkg/tracing"
	"net"
	"strconv"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	host           string
	tracer         opentracing.Tracer
	metricsFactory metrics.Factory
	logger         log.Factory
	server         *grpc.Server
}

func NewServer(configs *config.Publisher, metricsFactory metrics.Factory, logger log.Factory) *Server {
	host := net.JoinHostPort("0.0.0.0", strconv.Itoa(configs.Port))
	// create tracer
	tracer := tracing.Init(configs.ServiceName, metricsFactory, logger)
	// grpc server
	server := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)),
		grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(tracer)),
	)
	return &Server{
		host:           host,
		tracer:         tracer,
		logger:         logger,
		metricsFactory: metricsFactory,
		server:         server,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.host)
	if err != nil {
		s.logger.Bg().Error("Create tcp listener failed", zap.Error(err))
		return err
	}
	RegisterPublisherServiceServer(s.server, &PublisherServiceServerImpl{logger: s.logger})
	s.logger.Bg().Info("Starting grpc server", zap.String("host", s.host))
	return s.server.Serve(lis)
}

// PublisherServiceServerImpl
type PublisherServiceServerImpl struct {
	logger log.Factory
}

func (p *PublisherServiceServerImpl) Echo(ctx context.Context, req *EchoRequest) (*EchoResponse, error) {
	p.logger.For(ctx).Info("Starting echo process", zap.String("helloStr", req.HelloStr), zap.Int32("num", req.Num))
	helloArr := make([]string, req.Num)
	for i := 0; i < int(req.Num); i++ {
		helloArr[i] = req.HelloStr + " " + strconv.Itoa(i)
	}
	p.logger.For(ctx).Info("Echo response", zap.Strings("helloArr", helloArr))
	return &EchoResponse{HelloStr: helloArr}, nil
}
