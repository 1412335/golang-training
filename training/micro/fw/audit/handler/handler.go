package handler

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	microBroker "github.com/micro/micro/v3/service/broker"
	logger "github.com/micro/micro/v3/service/logger"

	pb "fw/audit/proto"
	"fw/pkg/broker"
	"gorm.io/gorm"
)

type Audit struct {
	ID           string
	Topic        string
	Service      string
	ActionType   string
	ActionFunc   string
	ActionTime   time.Time
	PerformedBy  string
	ObjectName   string
	Object       string
	Recordedtime time.Time
}

type auditHandler struct {
	db     *gorm.DB
	broker *broker.Broker
	topic  string
}

func New(db *gorm.DB, broker *broker.Broker) *auditHandler {
	return &auditHandler{
		db:     db,
		broker: broker,
	}
}

func (a *auditHandler) SubscribeMessage(topic, queueName string) error {
	if a.broker == nil {
		return errors.New("no broker provided")
	}
	a.topic = topic
	return a.broker.SubMsg(topic, a.processMessage, queueName)
}

func (a *auditHandler) processMessage(msg *microBroker.Message) error {
	logger.Info("[sub] received message:", string(msg.Body), "header", msg.Header)

	topic := a.topic
	header := msg.Header
	body := string(msg.Body)

	if err := a.createAuditRecord(topic, header, body); err != nil {
		return err
	}
	return nil
}

func (a *auditHandler) createAuditRecord(topic string, header map[string]string, body string) error {
	actionTime, err := time.Parse(time.RFC3339, header["actionTime"])
	if err != nil {
		return errors.New("Parse time failed")
	}
	rec := &Audit{
		ID:           uuid.New().String(),
		Topic:        topic,
		Service:      header["service"],
		ActionType:   header["actionType"],
		ActionFunc:   header["actionFunc"],
		ActionTime:   actionTime,
		PerformedBy:  header["performedBy"],
		ObjectName:   header["objectName"],
		Object:       body,
		Recordedtime: time.Now(),
	}
	if err := a.db.Create(rec).Error; err != nil {
		logger.Errorf("Error connecting from db: %v", err)
		// return ErrConnectDB
		return err
	}
	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *auditHandler) Call(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logger.Info("Received Audit.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}
