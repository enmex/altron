package config

import (
	"altron/pkg/sftp"
	"fmt"
	"os"
	"strings"
)

func NewSFTPConfig() *sftp.Config {
	return &sftp.Config{
		APIServer:     fmt.Sprintf("%s:80", os.Getenv("SFTP_SFTPGO_SERVICE_HOST")),
		Server:        fmt.Sprintf("%s:22", os.Getenv("SFTP_SFTPGO_SERVICE_HOST")),
		User:          strings.ReplaceAll(os.Getenv("HOSTNAME"), "-", ""),
		Password:      strings.ReplaceAll(os.Getenv("HOSTNAME"), "-", ""),
		AdminPassword: os.Getenv("SFTP_ADMIN_PASS"),
	}
}
