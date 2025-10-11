package dto

// Response represents a generic API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents error information in response
type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// SuccessResponse creates a success response
func SuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(code, message string, details interface{}) Response {
	return Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ListResponse represents a paginated list response
type ListResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}
