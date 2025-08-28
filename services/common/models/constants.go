package models

type WorkspaceStatus string

var (
	LISTENING WorkspaceStatus = "LISTENING"
	WAITING   WorkspaceStatus = "WAITING"
	COMPLETED WorkspaceStatus = "COMPLETED"
)

var ComponentNames = []string{"ttl", "requests", "ua", "timestamps"}
