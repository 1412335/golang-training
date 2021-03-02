package broker

import (
	"errors"
	microBroker "github.com/micro/micro/v3/service/broker"
	logger "github.com/micro/micro/v3/service/logger"
)

type Broker struct {
	br microBroker.Broker
}

func New(br microBroker.Broker) *Broker {
	return &Broker{br: br}
}

func (b *Broker) connect() error {
	if b.br == nil {
		err1 := errors.New("broker is not defined or connected")
		return err1
	}
	return b.br.Connect()
}

func (b *Broker) Disconnect() {
	err := b.br.Disconnect()
	if err != nil {
		logger.Infof("Disconnect broker error: %v", err)
	}
}

// SendMsg sends message to broker so that is can be picked up by a subscription at some point. This is setup to be fire and forget
func (b *Broker) SendMsg(objectToSend []byte, header map[string]string, topic string) error {
	if err := b.connect(); err != nil {
		return err
	}

	var message microBroker.Message
	message.Header = header
	message.Body = objectToSend

	err := b.br.Publish(topic, &message)
	if err != nil {
		return err
	}
	logger.Infof("sent message to Topic %s. Message %v", topic, &header)
	return nil
}

// subscribe a topic with queue specified
func (b *Broker) SubMsg(topic string, handler func(*microBroker.Message) error, queueName string) error {
	if err := b.connect(); err != nil {
		return err
	}
	var opts []microBroker.SubscribeOption
	if queueName != "" {
		opts = append(opts, microBroker.Queue(queueName))
	}
	_, err := b.br.Subscribe(topic, handler, opts...)
	return err
}
