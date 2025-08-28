package models

import "errors"

var (
	ErrorUserNotFound                    = errors.New("user not found")
	ErrorUserAlreadyExists               = errors.New("user already exists")
	ErrorInvalidCredentials              = errors.New("invalid credentials")
	ErrorServiceAlreadyExists            = errors.New("service with this port already exists")
	ErrorServiceNotFound                 = errors.New("service not found")
	ErrorWorkspaceAlreadyExists          = errors.New("workspace with this name already exists")
	ErrorWorkspaceNotFound               = errors.New("workspace not found")
	ErrorSessionNotFound                 = errors.New("session not found")
	ErrorDuplicateSession                = errors.New("this session already exists")
	ErrorAllFieldsEmpty                  = errors.New("one of the field must be not empty")
	ErrorUnableToPutPortOnSessionTree    = errors.New("unable to put port on session tree")
	ErrorUnableToDeletePortOnSessionTree = errors.New("unable to delete port on session tree")
	ErrorDuplicateFilter                 = errors.New("filter already exists")
	ErrorUnknownExportType               = errors.New("unknown export type")
	ErrorInvalidSearchInput              = errors.New("invalid search input")
	ErrorBlockingFilterServiceNotSet     = errors.New("service not set for blocking filter is empty")
	ErrorUnableToStopLogsListening       = errors.New("unable to stop listening to logs")
	ErrorNoFilesInPacket                 = errors.New("no files in packet")
	ErrorImportPcap                      = errors.New("something went wrong while importing pcap")
	ErrorSomethingWrong                  = errors.New("something went wrong")
	ErrorOnlyPcapSupported               = errors.New("only pcap files are supported")
	ErrorInvalidFilename                 = errors.New("invalid file name")
	ErrorPcapWorkspaceNotFound           = errors.New("pcap workspace not found")
)
