package auth

import (
	"context"
	repo "udo-golang/internal/adapters/mongo/repositories/user"

	jwt "github.com/dgrijalva/jwt-go"
)

type SignUpDto struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	IsAdmin   bool   `json:"isAdmin"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type ForgotPasswordDTO struct {
	Email string `json:"email" binding:"required"`
}
type ResetPasswordDTO struct {
	Password string `json:"password" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

type SignedDetails struct {
	Email   string
	ID      string
	IsAdmin bool
	jwt.StandardClaims
}

type Service interface {
	Signup(ctx context.Context, u *repo.User, body SignUpDto) (*repo.User, *string, error)
	Login(ctx context.Context, u LoginDTO) (*repo.User, *string, error)
	// GenerateToken(email, uid, subscription string) (signedToken string, signedRefreshToken string, err error)
	// ForgotPassword(ctx context.Context, req ForgotPasswordDTO) error
	// ResetPassword(ctx context.Context, req ResetPasswordDTO) error
}
