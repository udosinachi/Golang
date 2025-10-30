package user

import (
	"context"
	"errors"
	"os"
	repo "udo-golang/internal/adapters/mongo/repositories/user"
	commonErrors "udo-golang/internal/common/errors"

	"go.mongodb.org/mongo-driver/bson"
)

type server struct {
	userRepo repo.Repository
}

func NewUserService(userRepo repo.Repository) Server {
	return &server{userRepo}
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func (s *server) AllUsers(ctx context.Context, page, pageSize int) ([]repo.User, int64, error) {
	users, err := s.userRepo.GetAllUsersRepo(ctx, page, pageSize, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	count, err := s.userRepo.GetUserCountRepo(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return users, int64(count), nil
}

func (s *server) Create(c context.Context, u repo.User) (*repo.User, error) {
	_, err := s.userRepo.GetByEmailRepo(c, u.Email)

	if err != nil {
		return nil, err
	}

	if len(u.Password) < 8 {
		return nil, commonErrors.ErrShortPassword
	}

	user, err := s.userRepo.CreateUserRepo(c, u)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (s *server) GetByID(ctx context.Context, id string) (*repo.User, error) {
	if id == "" {
		return nil, errors.New("ID is required")
	}

	user, err := s.userRepo.GetUserByIDRepo(ctx, id)

	if err != nil {
		return nil, err
	}

	return user, err
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
