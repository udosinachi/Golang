package user

import (
	"context"
	"errors"
	"os"
	repo "udo-golang/internal/adapters/mongo/repositories/user"
	commonErrors "udo-golang/internal/common/errors"
)

type server struct {
	userRepo repo.Repository
}

func NewServer(userRepo repo.Repository) Server {
	return &server{userRepo}
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func (s *server) Create(c context.Context, u repo.User) (*repo.User, error) {
	_, err := s.userRepo.GetByEmail(c, u.Email)

	if err != nil {
		return nil, err
	}

	if len(u.Password) < 8 {
		return nil, commonErrors.ErrShortPassword
	}

	user, err := s.userRepo.Create(c, u)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (s *server) GetByID(ctx context.Context, id string) (*repo.User, error) {
	if id == "" {
		return nil, errors.New("ID is required")
	}

	user, err := s.userRepo.GetUserByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *server) List(ctx context.Context) ([]repo.User, int64, error) {
	var users = make([]repo.User, 5)

	return users, 7, nil
}

func (s *server) GetUser(ctx context.Context) (repo.User, error) {
	var user repo.User

	return user, nil
}

func (s *server) Update(ctx context.Context, t repo.User) (*repo.User, error) {
	var user repo.User

	return &user, nil
}

func (s *server) Delete(ctx context.Context, id string) error {
	return nil
}
