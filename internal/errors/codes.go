package errors

const (
	TopolithErrorUnknown TopolithErrorCode = 500
)

var topolithErrorDescriptions = map[TopolithErrorCode]string{
	TopolithErrorUnknown: "An unknown error occurred",
}
