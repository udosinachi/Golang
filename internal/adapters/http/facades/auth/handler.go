package auth

import (
	"time"
	"udo-golang/internal/adapters/http/common"
	repo "udo-golang/internal/adapters/mongo/repositories/user"
	authService "udo-golang/internal/services/auth"

	"github.com/gin-gonic/gin"
)

type authResponse struct {
	User  repo.User `json:"user"`
	Token string    `json:"token"`
}

func (f *Facade) Login(c *gin.Context) {

	var user authService.LoginDTO
	err := c.BindJSON(&user)

	if err != nil {
		common.SendBadRequest(c, err.Error())
		return
	}

	fetchedUser, token, err := f.auth.Login(c, user)
	if err != nil {
		common.SendBadRequest(c, err.Error())
		return
	}

	response := authResponse{
		User:  *fetchedUser,
		Token: *token,
	}

	common.SendOk(c, response, "Request Successful")

}

func (f *Facade) Signup(c *gin.Context) {

	var body authService.SignUpDto
	if err := c.BindJSON(&body); err != nil {
		common.SendBadRequest(c, err.Error())
		return
	}

	user := repo.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
		IsAdmin:   body.IsAdmin,
		CreatedAt: time.Now(),
	}

	fetchedUser, token, err := f.auth.Signup(c, &user, body)
	if err != nil {
		common.SendBadRequest(c, err.Error())
		return
	}

	response := authResponse{
		User:  *fetchedUser,
		Token: *token,
	}

	common.SendOk(c, response, "Request Successful")

}
