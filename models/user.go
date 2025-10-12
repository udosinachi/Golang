package models

import (
	"time"
	"udo-golang/database"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FirstName  string             `bson:"firstName" validate:"required"`
	LastName   string             `bson:"lastName"  validate:"required"`
	Email      string             `bson:"email" validate:"required,email"`
	Password   string             `bson:"password" validate:"required,min=6"`
	IsAdmin    bool               `bson:"isAdmin" default:"false"`
	IsVerified bool               `bson:"isVerified" default:"false"`
	LastLogin  *time.Time         `bson:"lastLogin"`
	Otp        *int               `bson:"otp"`
	OtpExpire  *time.Time         `bson:"otpExpire"`
	CreatedAt  time.Time          `bson:"createdAt"`
	UpdatedAt  *time.Time         `bson:"updatedAt"`
}

// create a validator instance
var validate = validator.New()

func (u *User) ValidateUser() error {
	return validate.Struct(u)
}

var UserCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
