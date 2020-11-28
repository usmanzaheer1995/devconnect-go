package models

import "strings"

const (
	// ErrNotFound is returned when a resource cannot be found
	// in the database
	ErrNotFound modelError = "models: resource not found"

	// ErrPasswordIncorrect is returned when an invalid password
	// is used when attempting to authenticate a user.
	ErrPasswordIncorrect modelError = "models: incorrect password provided"

	// ErrEmailRequired is returned when an email address is not
	// provided when creating a user
	ErrEmailRequired modelError = "models: Email address is required"

	// ErrEmailInvalid is returned when an email address provided
	// does not match our requirements
	ErrEmailInvalid modelError = "models: Email address is not valid"

	// ErrEmailTaken is returned when an update or create is attempted
	// with an email address that is already in use.
	ErrEmailTaken modelError = "models: email address is already taken"

	// ErrPasswordTooShort is returned when an update or create is attempted with a password
	// that is less than 8 characters
	ErrPasswordTooShort modelError = "models: password must be atleast 8 characters in length"

	// ErrPasswordRequired is returned when a create is attempted without a user
	// password provided
	ErrPasswordRequired modelError = "models: password is required"

	ErrTitleRequired modelError = "models: title is required"

	// ErrRememberTooShort is returned when a remember token is not
	// at least 32 bytes
	ErrRememberTooShort PrivateError = "models: remember token must be atleast 32 bytes"

	// ErrRememberRequired is returned when a create or update is attempted without a user
	// remember token hash is provided
	ErrRememberRequired PrivateError = "models: remember token is required"

	ErrUserIDRequired PrivateError = "models: user ID is required"

	// ErrIDInvalid is returned whan an invalid ID is provided
	// to a method like Delete
	ErrIDInvalid PrivateError = "models: ID must be > 0"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	return strings.Title(s)
}

type PrivateError string

func (e PrivateError) Error() string {
	return string(e)
}

type PublicError interface {
	error
	Public() string
}
