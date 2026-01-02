package doc

type Filter struct {
	// ColumnField is the name of the column
	ColumnField string `json:"column_field"`
	// OperatorValue is the operator to use for this filtering
	OperatorValue string `json:"operator_value"`
	// Value is the filter value for this filter
	Value string `json:"value"`
}
type Sort struct {
	// Field is the name of the field to sort
	Field string `json:"field"`
	// Sort is the type of sort to use
	// Sort order:
	// * asc - Ascending, from A to Z.
	// * desc - Descending, from Z to A.
	Sort string `json:"sort"`
}
type FilterParams struct {
	// Sort contains list of sorting methods
	Sort []Sort `json:"sort"`
	// Page specifies which page number to return
	Page int `json:"page"`
	// PerPage specifies number of results per page
	PerPage int `json:"per_page"`
	// Filter contains list of filters to apply
	Filter []Filter `json:"filter"`
	// Search specifies string value to search across all columns
	Search string `json:"search"`
	// LinkOperator is the operator used to connect multiple filters
	LinkOperator string `json:"link_operator"`
}

type PgnFltQueryParams struct {
	// Sort holds an object string of type Sort
	Sort string `json:"sort" form:"sort"`
	// Filter holds an array of object of type Filter
	Filter string `json:"filter" form:"filter"`
	// Search holds a search string
	Search string `json:"search" form:"search"`
	// Page is the page number to return
	Page int `json:"page" form:"page"`
	// PerPage is the number of values per page
	PerPage int `json:"per_page" form:"per_page"`
	// LinkOperator is the operator used to link multiple filters
	LinkOperator string `json:"link_operator" form:"link_operator"`
}

type SuccessResponse struct {
	// Success is only true if the request was successful.
	Success bool `json:"success" default:"false"`
	// MetaData contains additional data like filtering, pagination, etc.
	MetaData
	// Data contains the actual data of the response.
	Data interface{} `json:"data"`
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
	// Success is only true if the request was successful.
	Success bool `json:"success" default:"false"`
	// Error contains the error detail if the request was not successful.
	Error struct { // Code is the error code. It is not status code
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
	} `json:"error"`
}

type FieldError struct {
	// Name is the name of the field that caused the error.
	Name string `json:"name"`
	// Description is the error description for this field.
	Description string `json:"description"`
}
