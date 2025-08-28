package models

import "errors"

var (
	ErrorAddingSessionsToWorkspace   = errors.New("something went wrong when adding sessions to workspace")
	ErrorCreatingSessionAnalyzerData = errors.New("something went wrong when creating session analyzer data")
	ErrorSomethingWrong              = errors.New("something went wrong")
	ErrorAgentAlreadyConnected       = errors.New("agent already connected")
	ErrorAgentNotConnected           = errors.New("agent not connected")
)
