package user

import (
	"context"
	"errors"
	"os"
	"time"
	repo "udo-golang/internal/adapters/mongo/repositories/user"

	"go.mongodb.org/mongo-driver/bson"
)

type server struct {
	userRepo repo.Repository
}

func NewUserService(userRepo repo.Repository) Server {
	return &server{userRepo}
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func (s *server) AllUsers(ctx context.Context, page, pageSize int, filter bson.M) ([]repo.User, int64, error) {

	users, err := s.userRepo.GetAllUsersRepo(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.userRepo.GetUserCountRepo(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, int64(count), nil
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

func (s *server) Update(ctx context.Context, body UpdateUserDTO, id string) (*repo.User, error) {
	if id == "" {
		return nil, errors.New("ID is required")
	}
	update := bson.M{
		"firstName": body.FirstName,
		"lastName":  body.LastName,
		"isAdmin":   body.IsAdmin,
		"updatedAt": time.Now(),
	}

	updatedUser, err := s.userRepo.UpdateUserRepo(ctx, id, update)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *server) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("ID is required")
	}

	_, userErr := s.userRepo.GetUserByIDRepo(ctx, id)

	if userErr != nil {
		return userErr
	}

	err := s.userRepo.DeleteUserRepo(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
