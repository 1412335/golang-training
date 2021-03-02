package broker

import (
	"errors"
	"fmt"
	microBroker "github.com/micro/micro/v3/service/broker"
	logger "github.com/micro/micro/v3/service/logger"
)

type Broker struct {
	br microBroker.Broker
}

func New(br microBroker.Broker) *Broker {
	return &Broker{br: br}
}

func (b *Broker) Disconnect() {
	b.br.Disconnect()
}

// SendMsg sends message to broker so that is can be picked up by a subscription at some point. This is setup to be fire and forget
func (b *Broker) SendMsg(objectToSend []byte, header map[string]string, topic string) error {

	var message microBroker.Message
	message.Header = header
	message.Body = objectToSend

	if b.br == nil {
		err1 := errors.New("broker is not defined or connected")
		return err1
	}
	err := b.br.Connect()
	if err != nil {
		return err
	}
	err = b.br.Publish(topic, &message)
	if err != nil {
		return err
	}
	logger.Infof("sent message to Topic %s. Message %v", topic, &header)
	return nil
}

// subscribe a topic with queue specified
func (b *Broker) SubMsg(topic, queueName string) error {
	var opts []microBroker.SubscribeOption
	if queueName != "" {
		opts = append(opts, microBroker.Queue(queueName))
	}
	_, err := b.br.Subscribe(topic, func(p *microBroker.Message) error {
		fmt.Println("[sub] received message:", string(p.Body), "header", p.Header)
		return nil
	}, opts...)
	return err
}
