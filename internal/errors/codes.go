package errors

const (
	TopolithErrorUnknown  TopolithErrorCode = 500
	TopolithErrorInvalid                    = 400
	TopolithErrorNotFound                   = 404
	TopolithErrorConflict                   = 409
)

var topolithErrorDescriptions = map[TopolithErrorCode]string{
	TopolithErrorUnknown:  "An unknown error occurred",
	TopolithErrorInvalid:  "Invalid input",
	TopolithErrorNotFound: "Not found",
	TopolithErrorConflict: "Conflict or impossible state",
}
