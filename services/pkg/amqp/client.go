package amqp

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type Client struct {
	prefix string
	conn   *amqp.Connection
}

func NewClient(cfg *Config) (*Client, error) {
	amqpConn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", cfg.User, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}
	amqpConn.Config.Heartbeat = 5 * time.Second
	return &Client{
		prefix: cfg.Prefix,
		conn:   amqpConn,
	}, nil
}

func (c *Client) Channel() (*amqp.Channel, error) {
	return c.conn.Channel()
}

func (c *Client) Prefix() string {
	return c.prefix
}
