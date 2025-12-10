package utils

import (
	"github.com/gin-gonic/gin"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// PageMeta represents pagination metadata
type PageMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int64 `json:"totalPages"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessResponseWithMeta sends a successful response with pagination metadata
func SuccessResponseWithMeta(c *gin.Context, statusCode int, message string, data interface{}, meta interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, errors interface{}) {
	c.JSON(400, APIResponse{
		Success: false,
		Message: "Validation failed",
		Data:    errors,
	})
}

// UnauthorizedResponse sends an unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	ErrorResponse(c, 401, message)
}

// ForbiddenResponse sends a forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	ErrorResponse(c, 403, message)
}

// NotFoundResponse sends a not found response
func NotFoundResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	ErrorResponse(c, 404, message)
}

// ConflictResponse sends a conflict response
func ConflictResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Resource already exists"
	}
	ErrorResponse(c, 409, message)
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	ErrorResponse(c, 500, message)
}

// CalculateTotalPages calculates total pages for pagination
func CalculateTotalPages(total int64, limit int) int64 {
	if limit == 0 {
		return 0
	}
	pages := total / int64(limit)
	if total%int64(limit) > 0 {
		pages++
	}
	return pages
}

// NewPageMeta creates a new PageMeta instance
func NewPageMeta(total int64, page, limit int) PageMeta {
	return PageMeta{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: CalculateTotalPages(total, limit),
	}
}
