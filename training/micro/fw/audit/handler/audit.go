package handler

import (
	"context"

	log "github.com/micro/micro/v3/service/logger"

	audit "fw/audit/proto"
)

type Audit struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Audit) Call(ctx context.Context, req *audit.Request, rsp *audit.Response) error {
	log.Info("Received Audit.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Audit) Stream(ctx context.Context, req *audit.StreamingRequest, stream audit.Audit_StreamStream) error {
	log.Infof("Received Audit.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&audit.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Audit) PingPong(ctx context.Context, stream audit.Audit_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&audit.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
