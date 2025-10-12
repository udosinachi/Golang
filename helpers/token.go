package helpers

import (
	"fmt"
	"log"
	"udo-golang/database"

	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email        string
	ID           string
	Subscription string
	Type         string
	jwt.StandardClaims
}
type SignedEmployeeDetails struct {
	Email      string
	ID         string
	Type       string
	BusinessID string
	jwt.StandardClaims
}

type AuthService interface {
	// VerifyToken checks the provided token, returning the corresponding user
	// if valid or an error if not.
	ValidateToken(signedToken string) (claims *SignedDetails, msg string)
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("JWT_SECRET_KEY")

func GenerateAllTokens(email string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:          email,
		ID:             uid,
		Type:           "user",
		StandardClaims: jwt.StandardClaims{},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	// if claims.ExpiresAt < time.Now().Local().Unix(){
	// 	msg = fmt.Sprintf("token is expired")
	// 	msg = err.Error()
	// 	return
	// }
	return claims, msg
}
