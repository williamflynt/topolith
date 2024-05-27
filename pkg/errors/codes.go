package errors

const (
	TopolithErrorUnknown      TopolithErrorCode = 500
	TopolithErrorInvalid                        = 400
	TopolithErrorNotFound                       = 404
	TopolithErrorConflict                       = 409
	TopolithErrorBadSyncState                   = 502
	TopolithErrorMultiple                       = 600
)

var topolithErrorDescriptions = map[TopolithErrorCode]string{
	TopolithErrorUnknown:      "An unknown error occurred",
	TopolithErrorInvalid:      "Invalid input",
	TopolithErrorNotFound:     "Not found",
	TopolithErrorConflict:     "Conflict or impossible state",
	TopolithErrorBadSyncState: "Issue with World state sync detected",
	TopolithErrorMultiple:     "Multiple errors",
}
