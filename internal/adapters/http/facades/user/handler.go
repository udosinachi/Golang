package user

import (
	"strconv"
	"time"
	"udo-golang/internal/adapters/http/common"
	userService "udo-golang/internal/services/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (f *Facade) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	search := c.Query("search")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	isAdmin := c.Query("isAdmin")
	isVerified := c.Query("isVerified")

	filter := bson.M{}

	if search != "" {
		filter["$or"] = []bson.M{
			{"firstName": bson.M{"$regex": search, "$options": "i"}},
			{"lastName": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	dateFilter := bson.M{}
	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			dateFilter["$gte"] = start
		}
	}
	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			dateFilter["$lte"] = end.Add(24 * time.Hour)
		}
	}
	if len(dateFilter) > 0 {
		filter["createdAt"] = dateFilter
	}

	if isAdmin != "" {
		isAdminValue, err := strconv.ParseBool(isAdmin)
		if err == nil {
			filter["isAdmin"] = isAdminValue
		}
	}

	if isVerified != "" {
		isVerifiedValue, err := strconv.ParseBool(isVerified)
		if err == nil {
			filter["isVerified"] = isVerifiedValue

		}

	}

	users, total, err := f.service.AllUsers(c.Request.Context(), page, pageSize, filter)
	if err != nil {
		common.SendBadRequest(c, "Failed to fetch users")
		return
	}

	response := gin.H{
		"data": users,
		"pagination": gin.H{
			"page":        page,
			"per_page":    pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}

	common.SendOk(c, response, "Users fetched successfully")
}

func (f *Facade) GetUsersById(c *gin.Context) {
	id := c.Param("id")
	user, err := f.service.GetByID(c, id)

	if err != nil {
		common.SendBadRequest(c, err.Error())
		return
	}

	common.SendOk(c, user, "User Fetched Successfully")
}

func (f *Facade) DeleteUserById(c *gin.Context) {
	id := c.Param("id")
	err := f.service.Delete(c, id)

	if err != nil {
		common.SendBadRequest(c, err.Error())
		return
	}

	common.SendOk(c, []string{}, "User Deleted Successfully")
}

func (f *Facade) UpdateUserById(c *gin.Context) {
	id := c.Param("id")
	var body userService.UpdateUserDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		common.SendBadRequest(c, "Invalid request body")
		return
	}

	user, err := f.service.Update(c, body, id)

	if err != nil {
		common.SendBadRequest(c, err.Error())
		return
	}

	common.SendOk(c, user, "User Updated Successfully")
}
