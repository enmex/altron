package amqp

type Config struct {
	Host     string `yaml:"amqp-host"`
	Port     int    `yaml:"amqp-port"`
	User     string `yaml:"amqp-user"`
	Password string `yaml:"amqp-password"`
	Prefix   string
}
