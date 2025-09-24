package config

import (
	"altron/pkg/sftp"
	"fmt"
	"strings"
)

func NewSFTPConfig() *sftp.Config {
	return &sftp.Config{
		APIServer:     fmt.Sprintf("%s:8080", "altron.sftp.loc"), //os.Getenv("SFTP_SFTPGO_SERVICE_HOST")),
		Server:        fmt.Sprintf("%s:21", "altron.sftp.loc"), //os.Getenv("SFTP_SFTPGO_SERVICE_HOST")),
		User:          strings.ReplaceAll("altron", "-", ""),
		Password:      strings.ReplaceAll("altron", "-", ""),
		AdminPassword: "admin", //os.Getenv("SFTP_ADMIN_PASS"),
	}
}
