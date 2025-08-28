package docker

import "time"

type SendLogResponse struct {
	SentAt  time.Time `json:"sentAt"`
	Message string    `json:"message"`
}
