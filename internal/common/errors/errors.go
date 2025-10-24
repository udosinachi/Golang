package errors

import "errors"

var (
	ErrInvalidToken      = errors.New("invalid access token")
	ErrAccessDenied      = errors.New("access denied")
	ErrIncorrectPassword = errors.New("the email/password you entered is incorrect, please try again")
	ErrUnauthorized      = errors.New("user unauthorized")
	ErrGenericError      = errors.New("an error occured")
	ErrInternal          = errors.New("internal server error occured")
	ErrShortPassword     = errors.New("password is short")
	ErrWrongEmail        = errors.New("the email you entered is incorrect")
)
