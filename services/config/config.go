package config

import (
	"altron/pkg/amqp"
	"altron/pkg/auth"
	"altron/pkg/db"
	"altron/pkg/plugin"
	"altron/pkg/redis"
	"altron/pkg/sftp"
)

type Config struct {
	App    *AppConfig
	Auth   *auth.Config
	DB     *db.Config
	Plugin *plugin.Config
	AMQP   *amqp.Config
	SFTP   *sftp.Config
	Redis  *redis.Config
}

func NewConfig(
	app *AppConfig,
	auth *auth.Config,
	db *db.Config,
	plugin *plugin.Config,
	amqp *amqp.Config,
	sftp *sftp.Config,
	redis *redis.Config,
) *Config {
	return &Config{
		App:    app,
		Auth:   auth,
		DB:     db,
		Plugin: plugin,
		AMQP:   amqp,
		SFTP:   sftp,
		Redis:  redis,
	}
}
