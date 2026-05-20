package response

import (
	"time"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool         `json:"success"`
	Error   *ErrorDetail `json:"error"`
}

// ErrorDetail contains detailed error information
type ErrorDetail struct {
	Code      string          `json:"code"`
	Message   string          `json:"message"`
	Details   []ErrorIssue    `json:"details,omitempty"`
	TraceID   string          `json:"traceId"`
	Timestamp string          `json:"timestamp"`
}

// ErrorIssue represents a specific error issue
type ErrorIssue struct {
	Service string `json:"service"`
	Issue   string `json:"issue"`
}

// Pagination represents pagination information
type Pagination struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

// PaginatedData represents paginated list data
type PaginatedData struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

// AsyncJobData represents async job response data
type AsyncJobData struct {
	JobID            string `json:"job_id"`
	Status           string `json:"status"`
	EstimatedSeconds int    `json:"estimated_seconds,omitempty"`
}

// Success creates a successful response with data
func Success(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Data:    data,
	}
}

// SuccessList creates a successful paginated list response
func SuccessList(items interface{}, pagination Pagination) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Data: PaginatedData{
			Items:      items,
			Pagination: pagination,
		},
	}
}

// SuccessAsync creates a successful async job response
func SuccessAsync(jobID, status string, estimatedSeconds int) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Data: AsyncJobData{
			JobID:            jobID,
			Status:           status,
			EstimatedSeconds: estimatedSeconds,
		},
	}
}

// Error creates an error response
func Error(code, message, traceID string, details ...ErrorIssue) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Error: &ErrorDetail{
			Code:      code,
			Message:   message,
			Details:   details,
			TraceID:   traceID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

// NewPagination creates a pagination object with calculated total pages
func NewPagination(total, page, pageSize int) Pagination {
	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	return Pagination{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// Common Error Codes
const (
	// Zone B - Web Application (BFF)
	BNYBBadRequest          = "B-NYB-400"
	BNYBUnauthorized        = "B-NYB-401"
	BNYBForbidden           = "B-NYB-403"
	BNYBNotFound            = "B-NYB-404"
	BNYBConflict            = "B-NYB-409"
	BNYBValidationFailed    = "B-NYB-422"
	BNYBTooManyRequests     = "B-NYB-429"
	BNYBInternalServerError = "B-NYB-500"
	BNYBServiceUnavailable  = "B-NYB-503"
	BNYBGatewayTimeout      = "B-NYB-504"

	// Zone C - Core Services
	// C-AUT - Auth Service
	CAUTBadRequest          = "C-AUT-400"
	CAUTUnauthorized        = "C-AUT-401"
	CAUTForbidden           = "C-AUT-403"
	CAUTNotFound            = "C-AUT-404"
	CAUTConflict            = "C-AUT-409"
	CAUTValidationFailed    = "C-AUT-422"
	CAUTTooManyRequests     = "C-AUT-429"
	CAUTInternalServerError = "C-AUT-500"
	CAUTServiceUnavailable  = "C-AUT-503"

	// C-LST - Listing Service
	CLSTBadRequest          = "C-LST-400"
	CLSTUnauthorized        = "C-LST-401"
	CLSTForbidden           = "C-LST-403"
	CLSTNotFound            = "C-LST-404"
	CLSTConflict            = "C-LST-409"
	CLSTValidationFailed    = "C-LST-422"
	CLSTTooManyRequests     = "C-LST-429"
	CLSTInternalServerError = "C-LST-500"
	CLSTServiceUnavailable  = "C-LST-503"

	// C-MED - Media Service
	CMEDBadRequest          = "C-MED-400"
	CMEDUnauthorized        = "C-MED-401"
	CMEDForbidden           = "C-MED-403"
	CMEDNotFound            = "C-MED-404"
	CMEDConflict            = "C-MED-409"
	CMEDPayloadTooLarge     = "C-MED-413"
	CMEDValidationFailed    = "C-MED-422"
	CMEDTooManyRequests     = "C-MED-429"
	CMEDInternalServerError = "C-MED-500"
	CMEDServiceUnavailable  = "C-MED-503"

	// C-SRC - Search Service
	CSRCBadRequest          = "C-SRC-400"
	CSRCUnauthorized        = "C-SRC-401"
	CSRCForbidden           = "C-SRC-403"
	CSRCNotFound            = "C-SRC-404"
	CSRCValidationFailed    = "C-SRC-422"
	CSRCTooManyRequests     = "C-SRC-429"
	CSRCInternalServerError = "C-SRC-500"
	CSRCServiceUnavailable  = "C-SRC-503"

	// C-NOT - Notification Service
	CNOTBadRequest          = "C-NOT-400"
	CNOTUnauthorized        = "C-NOT-401"
	CNOTForbidden           = "C-NOT-403"
	CNOTNotFound            = "C-NOT-404"
	CNOTValidationFailed    = "C-NOT-422"
	CNOTTooManyRequests     = "C-NOT-429"
	CNOTInternalServerError = "C-NOT-500"
	CNOTServiceUnavailable  = "C-NOT-503"

	// C-RPT - Report Service
	CRPTBadRequest          = "C-RPT-400"
	CRPTUnauthorized        = "C-RPT-401"
	CRPTForbidden           = "C-RPT-403"
	CRPTNotFound            = "C-RPT-404"
	CRPTValidationFailed    = "C-RPT-422"
	CRPTTooManyRequests     = "C-RPT-429"
	CRPTInternalServerError = "C-RPT-500"
	CRPTServiceUnavailable  = "C-RPT-503"

	// C-PRD - Product Service
	CPRDBadRequest          = "C-PRD-400"
	CPRDUnauthorized        = "C-PRD-401"
	CPRDForbidden           = "C-PRD-403"
	CPRDNotFound            = "C-PRD-404"
	CPRDValidationFailed    = "C-PRD-422"
	CPRDTooManyRequests     = "C-PRD-429"
	CPRDInternalServerError = "C-PRD-500"
	CPRDServiceUnavailable  = "C-PRD-503"

	// Zone X - Integration Service
	XINTBadRequest          = "X-INT-400"
	XINTUnauthorized        = "X-INT-401"
	XINTForbidden           = "X-INT-403"
	XINTNotFound            = "X-INT-404"
	XINTValidationFailed    = "X-INT-422"
	XINTTooManyRequests     = "X-INT-429"
	XINTInternalServerError = "X-INT-500"
	XINTBadGateway          = "X-INT-502"
	XINTServiceUnavailable  = "X-INT-503"
	XINTGatewayTimeout      = "X-INT-504"
)

// Helper functions for common errors

// BadRequest creates a 400 Bad Request error
func BadRequest(code, message, traceID string, details ...ErrorIssue) *ErrorResponse {
	return Error(code, message, traceID, details...)
}

// Unauthorized creates a 401 Unauthorized error
func Unauthorized(code, message, traceID string, details ...ErrorIssue) *ErrorResponse {
	return Error(code, message, traceID, details...)
}

// Forbidden creates a 403 Forbidden error
func Forbidden(code, message, traceID string, details ...ErrorIssue) *ErrorResponse {
	return Error(code, message, traceID, details...)
}

// NotFound creates a 404 Not Found error
func NotFound(code, message, traceID string, details ...ErrorIssue) *ErrorResponse {
	return Error(code, message, traceID, details...)
}

// Conflict creates a 409 Conflict error
func Conflict(code, message, traceID string, details ...ErrorIssue) *ErrorResponse {
	return Error(code, message, traceID, details...)
}

// PayloadTooLarge creates a 413 Payload Too Large error
func PayloadTooLarge(code, message, traceID string, details ...ErrorIssue) *ErrorResponse {
	return Error(code, message, traceID, details...)
}

// ValidationFailed creates a 422 Validation Failed error
func ValidationFailed(code, message, traceID string, details ...ErrorIssue) *ErrorResponse {
	return Error(code, message, traceID, details...)
}

// TooManyRequests creates a 429 Too Many Requests error
func TooManyRequests(code, message, traceID string, details ...ErrorIssue) *ErrorResponse {
	return Error(code, message, traceID, details...)
}

// InternalServerError creates a 500 Internal Server Error
func InternalServerError(code, message, traceID string, details ...ErrorIssue) *ErrorResponse {
	return Error(code, message, traceID, details...)
}

// BadGateway creates a 502 Bad Gateway error
func BadGateway(code, message, traceID string, details ...ErrorIssue) *ErrorResponse {
	return Error(code, message, traceID, details...)
}

// ServiceUnavailable creates a 503 Service Unavailable error
func ServiceUnavailable(code, message, traceID string, details ...ErrorIssue) *ErrorResponse {
	return Error(code, message, traceID, details...)
}

// GatewayTimeout creates a 504 Gateway Timeout error
func GatewayTimeout(code, message, traceID string, details ...ErrorIssue) *ErrorResponse {
	return Error(code, message, traceID, details...)
}
