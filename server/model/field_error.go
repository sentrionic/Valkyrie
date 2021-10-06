package model

// ErrorsResponse contains a list of FieldError
type ErrorsResponse struct {
	Errors []FieldError `json:"errors"`
} //@name ErrorsResponse

// FieldError is used to help extract validation errors
type FieldError struct {
	// The property containing the error
	Field string `json:"field"`
	// The specific error message
	Message string `json:"message"`
} //@name FieldError

// ErrorResponse holds a custom error for the application
type ErrorResponse struct {
	Error HttpError `json:"error"`
} //@name ErrorResponse

// HttpError returns the Http error type and the specific message
type HttpError struct {
	// The Http Response as a string
	Type string `json:"type"`
	// The specific error message
	Message string `json:"message"`
} //@name HttpError
