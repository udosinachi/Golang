package user

import (
	"context"

	repo "udo-golang/internal/adapters/mongo/repositories/user"

	"go.mongodb.org/mongo-driver/bson"
)

type UpdateUserDTO struct {
	FirstName string `json:"firstName" validate:"required,min=3"`
	LastName  string `json:"lastName" validate:"required,min=3"`
	IsAdmin   bool   `json:"isAdmin"`
}
type Server interface {
	GetByID(ctx context.Context, id string) (*repo.User, error)
	AllUsers(ctx context.Context, page, pageSize int, filter bson.M) ([]repo.User, int64, error)
	GetUser(ctx context.Context) (repo.User, error)
	Update(ctx context.Context, body UpdateUserDTO, id string) (*repo.User, error)
	Delete(ctx context.Context, id string) error
}
