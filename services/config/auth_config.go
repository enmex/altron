package config

import (
	"altron/pkg/auth"
	"altron/utils"
)

func NewAuthConfig() *auth.Config {
	return &auth.Config{
		Secret:   []byte(utils.RandomString(20)),
		HashSalt: "Aeg6MVjyNu2ANY3HAHVP0gD8EpUOQS9i",
	}
}
