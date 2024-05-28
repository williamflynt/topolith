package errors

const (
	TopolithErrorInvalid      TopolithErrorCode = 400
	TopolithErrorNotFound                       = 404
	TopolithErrorConflict                       = 409
	TopolithErrorCommandErr                     = 450
	TopolithErrorInternal                       = 500
	TopolithErrorBadSyncState                   = 502
	TopolithErrorMultiple                       = 600
)

var topolithErrorDescriptions = map[TopolithErrorCode]string{
	TopolithErrorInternal:     "An unknown error occurred",
	TopolithErrorInvalid:      "Invalid input",
	TopolithErrorNotFound:     "Not found",
	TopolithErrorConflict:     "Conflict or impossible state",
	TopolithErrorBadSyncState: "Issue with World state sync detected",
	TopolithErrorMultiple:     "Multiple errors",
	TopolithErrorCommandErr:   "Error while executing command",
}
