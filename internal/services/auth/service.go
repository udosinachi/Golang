package auth

import (
	"context"
	"errors"
	"strings"
	repo "udo-golang/internal/adapters/mongo/repositories/user"
	commonErrors "udo-golang/internal/common/errors"
	"udo-golang/internal/common/password"
	"udo-golang/internal/common/token"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	userRepo  repo.Repository
	secretKey string
}

func NewService(userRepo repo.Repository, secretKey string) Service {
	return &service{
		userRepo,
		secretKey,
	}
}

func (s *service) Signup(ctx context.Context, u *repo.User, body SignUpDto) (*repo.User, *string, error) {

	validationErr := u.ValidateUser()

	if validationErr != nil {
		return nil, nil, validationErr
	}

	u.Email = strings.ToLower(u.Email)

	existingUser, err := s.userRepo.GetByEmailRepo(ctx, u.Email)
	if err == nil && existingUser != nil {
		return nil, nil, errors.New("email already in use")
	}

	u.ID = primitive.NewObjectID()
	u.FirstName = body.FirstName
	u.LastName = body.LastName
	u.Email = strings.ToLower(body.Email)
	u.IsVerified = true
	u.IsAdmin = body.IsAdmin
	u.CreatedAt = time.Now()

	hashedPass, hashErr := password.HashPassword(u.Password)
	if hashErr != nil {
		return nil, nil, errors.New("unable to hash password")
	}
	u.Password = hashedPass

	sToken, _, err := token.GenerateAllTokens(u.Email, u.ID.Hex(), u.IsAdmin)
	if err != nil {
		return nil, nil, err
	}

	user, createErr := s.userRepo.CreateUserRepo(ctx, *u)
	if createErr != nil {
		return nil, nil, createErr
	}

	return user, &sToken, nil
}

func (s *service) Login(ctx context.Context, u LoginDTO) (*repo.User, *string, error) {
	u.Email = strings.ToLower(u.Email)

	user, err := s.userRepo.GetByEmailRepo(ctx, u.Email)

	if err != nil {
		return nil, nil, commonErrors.ErrWrongEmail
	}

	if user.Email == "" {
		return nil, nil, commonErrors.ErrWrongEmail
	}

	lastLogin := time.Now()
	user.LastLogin = &lastLogin

	check, msg := password.VerifyPassword(u.Password, user.Password)

	if !check {
		return nil, nil, errors.New(msg)
	}

	sToken, _, err := token.GenerateAllTokens(user.Email, user.ID.Hex(), user.IsAdmin)

	if err != nil {
		return nil, nil, err
	}

	return user, &sToken, nil

}
