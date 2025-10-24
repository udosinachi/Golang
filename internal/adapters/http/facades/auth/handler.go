package auth

import (
	"udo-golang/internal/adapters/http/common"
	repo "udo-golang/internal/adapters/mongo/repositories/user"
	authService "udo-golang/internal/services/auth"

	"github.com/gin-gonic/gin"
)

type authResponse struct {
	user  repo.User
	token string
}

func (f *Facade) Login(c *gin.Context) {

	var user authService.LoginDTO
	err := c.BindJSON(&user)

	if err != nil {
		common.SendNotFound(c, err.Error())
		return
	}

	fetchedUser, token, err := f.auth.Login(c, user)
	if err != nil {
		common.SendNotFound(c, err.Error())
		return
	}

	response := authResponse{
		user:  *fetchedUser,
		token: *token,
	}

	common.SendOk(c, response, "Request Successful")

}

func (f *Facade) Signup(c *gin.Context) {

	var user repo.User
	err := c.BindJSON(&user)

	if err != nil {
		common.SendNotFound(c, err.Error())
		return
	}

	bodySignup := authService.SignUpDto{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsAdmin:   user.IsAdmin,
		Password:  user.Password,
	}

	fetchedUser, token, err := f.auth.Signup(c, &user, bodySignup)
	if err != nil {
		common.SendNotFound(c, err.Error())
		return
	}

	response := authResponse{
		user:  *fetchedUser,
		token: *token,
	}

	common.SendOk(c, response, "Request Successful")

}
