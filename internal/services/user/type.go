package user

import (
	"context"

	repo "udo-golang/internal/adapters/mongo/repositories/user"
)

// type CreateUserDTO struct {
// 	ID               primitive.ObjectID `bson:"_id,omitempty"`
// 	FirstName        string             `bson:"first_name" validate:"required"`
// 	LastName         string             `bson:"last_name" validate:"required"`
// 	Email            string             `bson:"email" validate:"required,email"`
// 	Country          string             `bson:"country" validate:"required"`
// 	Phone            string             `bson:"phone" validate:"required"`
// 	Password         string             `bson:"password" validate:"required,min=6"`
// 	Role             string             `bson:"role" validate:"required,oneof=owner user"`
// 	SubscriptionTier string             `bson:"subscription_tier" validate:"required,oneof=Free Basic Standard Premium"`
// 	CreatedAt        time.Time          `bson:"created_at"`
// 	LastLoggedIn     time.Time          `bson:"last_logged_in"`
// }

type Server interface {
	Create(ctx context.Context, u repo.User) (*repo.User, error)
	GetByID(ctx context.Context, id string) (*repo.User, error)
	AllUsers(ctx context.Context, page, pageSize int) ([]repo.User, int64, error)
	GetUser(ctx context.Context) (repo.User, error)
	Update(ctx context.Context, t repo.User) (*repo.User, error)
	Delete(ctx context.Context, id string) error
}
