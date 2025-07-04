package response

import (
	"net/http"
	"time"

	"github.com/geoo115/property-manager/validator"
	"github.com/gin-gonic/gin"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination Pagination  `json:"pagination"`
	Timestamp  time.Time   `json:"timestamp"`
}

// Pagination contains pagination information
type Pagination struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	TotalItems int  `json:"total_items"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Error     interface{} `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Success response functions
func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func NoContent(c *gin.Context, message string) {
	c.JSON(http.StatusNoContent, APIResponse{
		Success:   true,
		Message:   message,
		Timestamp: time.Now(),
	})
}

// Error response functions
func BadRequest(c *gin.Context, message string, err interface{}) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success:   false,
		Message:   message,
		Error:     err,
		Timestamp: time.Now(),
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	})
}

func Conflict(c *gin.Context, message string, err interface{}) {
	c.JSON(http.StatusConflict, ErrorResponse{
		Success:   false,
		Message:   message,
		Error:     err,
		Timestamp: time.Now(),
	})
}

func UnprocessableEntity(c *gin.Context, message string, err interface{}) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		Success:   false,
		Message:   message,
		Error:     err,
		Timestamp: time.Now(),
	})
}

func TooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, ErrorResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	})
}

func InternalServerError(c *gin.Context, message string, err interface{}) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success:   false,
		Message:   message,
		Error:     err,
		Timestamp: time.Now(),
	})
}

// ValidationError handles validation errors
func ValidationError(c *gin.Context, errors validator.ValidationErrors) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success:   false,
		Message:   "Validation failed",
		Error:     errors,
		Timestamp: time.Now(),
	})
}

// Paginated response function
func Paginated(c *gin.Context, data interface{}, pagination Pagination, message string) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
		Timestamp:  time.Now(),
	})
}

// CalculatePagination calculates pagination values
func CalculatePagination(page, pageSize, totalItems int) Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	totalPages := (totalItems + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	return Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
