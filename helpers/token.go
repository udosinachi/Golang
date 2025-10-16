package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	Email   string `json:"email"`
	ID      string `json:"id"`
	IsAdmin bool   `json:"isAdmin"`
	jwt.StandardClaims
}
type AuthService interface {
	ValidateToken(signedToken string) (*SignedDetails, string)
}

var SECRET_KEY = os.Getenv("JWT_SECRET_KEY")

func GenerateAllTokens(email string, uid string, isAdmin bool) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:   email,
		ID:      uid,
		IsAdmin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// refresh token (expires in 3 days)
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Printf("Failed to sign access token: %v", err)
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Printf("Failed to sign refresh token: %v", err)
		return "", "", err
	}

	return token, refreshToken, nil
}

func ValidateToken(signedToken string) (*SignedDetails, string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, "Token has expired"
			}
			return nil, fmt.Sprintf("invalid token: %v", ve.Inner)
		}
		return nil, fmt.Sprintf("token parsing error: %v", err)
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok || !token.Valid {
		return nil, "The token is invalid"
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, "Token has expired"
	}

	return claims, ""
}
