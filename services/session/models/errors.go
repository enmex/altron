package models

import "errors"

var (
	ErrorServerNotFound              = errors.New("server with this port not found")
	ErrorClientNotFound              = errors.New("client with this host not found")
	ErrorAddingSessionsToWorkspace   = errors.New("something went wrong when adding sessions to workspace")
	ErrorCreatingSessionAnalyzerData = errors.New("something went wrong when creating session analyzer data")
	ErrorInterfaceNotFound           = errors.New("interface with this index not found")
	ErrorCachedLogsNotFound          = errors.New("cannot find cached logs for this container")
)
