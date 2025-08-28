package amqp

import (
	"sync"

	"github.com/streadway/amqp"
)

type Producer struct {
	prefix  string
	channel *amqp.Channel
	mut     *sync.Mutex
}

func NewProducer(client *Client) (*Producer, error) {
	ch, err := client.Channel()
	if err != nil {
		return nil, err
	}
	return &Producer{
		prefix:  client.Prefix(),
		channel: ch,
		mut:     &sync.Mutex{},
	}, nil
}

func (p *Producer) SendMessage(exchangeName string, message []byte) error {
	p.mut.Lock()
	defer p.mut.Unlock()
	return p.channel.Publish(
		p.prefix+exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (p *Producer) SendMessageToQueue(queueName string, message []byte) error {
	p.mut.Lock()
	defer p.mut.Unlock()
	return p.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}

func (p *Producer) CreateExchange(exchangeName string) error {
	p.mut.Lock()
	defer p.mut.Unlock()
	return p.channel.ExchangeDeclare(
		p.prefix+exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (p *Producer) DeleteExchange(exchangeName string) error {
	p.mut.Lock()
	defer p.mut.Unlock()
	return p.channel.ExchangeDelete(
		p.prefix+exchangeName,
		false,
		false,
	)
}

func (p *Producer) CreateQueue(queueName string) error {
	p.mut.Lock()
	defer p.mut.Unlock()
	_, err := p.channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (p *Producer) BindQueue(queueName, exchangeName string) error {
	p.mut.Lock()
	defer p.mut.Unlock()
	return p.channel.QueueBind(
		queueName,
		"",
		p.prefix+exchangeName,
		false,
		nil,
	)
}

func (p *Producer) DeleteQueue(queueName string) error {
	p.mut.Lock()
	defer p.mut.Unlock()
	_, err := p.channel.QueueDelete(
		queueName,
		false,
		false,
		false,
	)
	return err
}

func (p *Producer) Close() error {
	return p.channel.Close()
}
