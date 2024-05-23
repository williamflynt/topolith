package errors

import "fmt"

// TopolithErrorCode is an iota that represents the error code of a TopolithError.
type TopolithErrorCode int

// KvPair is a key-value pair where the Value can be represented as a string.
type KvPair struct {
	Key   string
	Value string
}

// TopolithError is an error type that is used in the Topolith Core, and implements the error interface.
// It contains a Code, Description, Message, and Data.
// The Message field is what we expect from basic Go errors (ex: errors.New("message")).
type TopolithError struct {
	Code        TopolithErrorCode // Code is the error code.
	Description string            // Description is a human-readable description of the error. It is generally tightly bound to the Code.
	Message     string            // Message is a human-readable message that is generally more detailed than Description.
	Data        []KvPair          // Data is a list of key-value pairs that provide additional context to the error.

	errs []error // errs is a list of errors that are wrapped by this error.
}

func (e TopolithError) UseCode(code TopolithErrorCode) TopolithError {
	e.Code = code
	if desc, ok := topolithErrorDescriptions[code]; ok {
		e.Description = desc
	}
	return e
}

func (e TopolithError) WithDescription(description string) TopolithError {
	e.Description = description
	return e
}

func (e TopolithError) WithData(data ...KvPair) TopolithError {
	e.Data = append(e.Data, data...)
	return e
}

func (e TopolithError) ClearData() TopolithError {
	e.Data = make([]KvPair, 0)
	return e
}

func (e TopolithError) WithError(err error) TopolithError {
	e.errs = append(e.errs, err)
	return e
}

// --- ERROR IMPLEMENTATION ---

// New returns a new TopolithError with the given text.
// It is designed to be used in the same way as errors.New.
func New(text string) TopolithError {
	return TopolithError{
		Code:        TopolithErrorUnknown,
		Description: topolithErrorDescriptions[TopolithErrorUnknown],
		Message:     text,
		Data:        make([]KvPair, 0),
		errs:        make([]error, 0),
	}
}

// Error returns a string representation of the TopolithError.
func (e TopolithError) Error() string {
	return fmt.Sprintf("TopolithError %d %s: %s", e.Code, e.Description, e.Message)
}

// Unwrap returns the first wrapped error, if any.
func (e TopolithError) Unwrap() error {
	if len(e.errs) == 0 {
		return nil
	}
	return e.errs[0]
}
