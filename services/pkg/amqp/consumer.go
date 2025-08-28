package amqp

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Consumer struct {
	prefix    string
	queueName string
	channel   *amqp.Channel
}

func NewConsumer(client *Client, queueName string) (*Consumer, error) {
	channel, err := client.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		prefix:    client.Prefix(),
		channel:   channel,
		queueName: queueName,
	}, nil
}

func (c *Consumer) Messages() (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		c.queueName,
		fmt.Sprintf("consumer-%s", c.queueName),
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *Consumer) Bind(exchangeName string) error {
	q, err := c.channel.QueueDeclare(
		c.queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return c.channel.QueueBind(
		q.Name,
		"",
		c.prefix+exchangeName,
		false,
		nil,
	)
}

func (c *Consumer) Close() error {
	_, err := c.channel.QueueDelete(
		c.queueName,
		false,
		false,
		false,
	)
	if err != nil {
		return err
	}
	return c.channel.Close()
}

func (c *Consumer) MemberID() string {
	return fmt.Sprintf("member-%s", c.queueName)
}
