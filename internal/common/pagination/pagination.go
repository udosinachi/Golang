package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pagination struct for handling paginated responses
type Pagination struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ExtractPagination extracts pagination parameters from the request
func ExtractPagination(c *gin.Context, defaultPerPage int) (int, int) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", strconv.Itoa(defaultPerPage)))
	if err != nil || perPage < 1 {
		perPage = defaultPerPage
	}

	return page, perPage
}

// CreatePaginationResponse builds pagination metadata
func CreatePaginationResponse(page, perPage int, total int64) Pagination {
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return Pagination{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}
