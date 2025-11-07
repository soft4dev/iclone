package errors

type messageType string

const (
	ErrorTypeError   = "error"
	ErrorTypeWarning = "warning"
	ErrorTypeInfo    = "info"
)

type CustomError struct {
	message     string
	ShowUsage   bool
	MessageType messageType
}

func NewCustomError(msg string, errorType messageType, showUsage bool) *CustomError {
	return &CustomError{
		message:     msg,
		MessageType: errorType,
		ShowUsage:   showUsage,
	}
}

func (e *CustomError) Error() string {
	return e.message
}
