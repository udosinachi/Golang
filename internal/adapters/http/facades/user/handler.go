package user

import (
	"strconv"
	"udo-golang/internal/adapters/http/common"

	"github.com/gin-gonic/gin"
)

func (f *Facade) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Add validation for page and pageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := f.service.AllUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		common.SendBadRequest(c, "Failed to fetch users")
		return
	}

	response := gin.H{
		"data": users,
		"pagination": gin.H{
			"currentPage": page,
			"pageSize":    pageSize,
			"total":       total,
			"pages":       (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}

	common.SendOk(c, response, "Users fetched successfully")
}
