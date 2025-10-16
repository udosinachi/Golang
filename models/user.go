package models

import (
	"time"
	"udo-golang/database"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName  string             `bson:"firstName" json:"firstName" validate:"required"`
	LastName   string             `bson:"lastName" json:"lastName" validate:"required"`
	Email      string             `bson:"email" json:"email" validate:"required,email"`
	Password   string             `bson:"password,omitempty" json:"-" validate:"required,min=6"`
	IsAdmin    bool               `bson:"isAdmin" json:"isAdmin"`
	IsVerified bool               `bson:"isVerified" json:"isVerified"`
	LastLogin  *time.Time         `bson:"lastLogin,omitempty" json:"lastLogin"`
	Otp        *string            `bson:"otp,omitempty" json:"otp"`
	OtpExpire  *time.Time         `bson:"otpExpire,omitempty" json:"otpExpire"`
	CreatedAt  time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt  *time.Time         `bson:"updatedAt,omitempty" json:"updatedAt"`
}

var validate = validator.New()

func (u *User) ValidateUser() error {
	return validate.Struct(u)
}

var UserCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
