package response

type Response struct {
	// OK is only true if the request was successful.
	Success bool `json:"success"`
	// MetaData contains additional data like filtering, pagination, etc.
	*MetaData
	// Data contains the actual data of the response.
	Data interface{} `json:"data"`
	// Error contains the error detail if the request was not successful.
	Error *ErrorResponse `json:"error"`
}

type MetaData struct {
	// Total number of items
	Count int `json:"count"`
	// URL for the next page, null if no next page
	Next *int `json:"next"`
	// URL for the previous page, null if no previous page
	Previous *int `json:"previous"`
}

type ErrorResponse struct {
	// Code is the error code. It is not status code
	Code int `json:"code"`
	// Message is the error message.
	Message string `json:"message"`
	// Description is the error description.
	Description string `json:"description"`
	// StackTrace is the stack trace of the error.
	// It is only returned for debugging
	StackTrace string `json:"stack_trace"`
	// FieldError is the error detail for each field, if available that is.
	FieldError []FieldError `json:"field_error"`
}

type FieldError struct {
	// Name is the name of the field that caused the error.
	Name string `json:"name"`
	// Description is the error description for this field.
	Description string `json:"description"`
}
