package user

import (
	errors2 "github.com/usmanzaheer1995/devconnect-go-v2/internal/errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/drexedam/gravatar"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userValidator struct {
	Errors map[string][]string
	UserDB
}

var _ UserDB = &userValidator{}

func newUserValidator(udb UserDB) *userValidator {
	return &userValidator{
		UserDB: udb,
		Errors: make(map[string][]string),
	}
}

type userValFunc func(*User) error

func runUserValFuncsArray(user *User, fns ...userValFunc) []error {
	var errs []error
	for _, fn := range fns {
		if err := fn(user); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func (uv *userValidator) setAvatar(u *User) error {
	url := gravatar.New(u.Email).
		Size(200).
		Default(gravatar.MysteryMan).
		Rating(gravatar.Pg).
		AvatarURL()
	u.Avatar = url
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		uv.Errors["email"] = append(uv.Errors["email"], "email is required")
		//return errors.New("email is required")
	}
	return nil
}

func (uv *userValidator) emailExists(user *User) error {
	existing, err := uv.UserDB.ByEmail(user.Email)
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if user.ID != existing.ID {
		uv.Errors["email"] = append(uv.Errors["email"], "email already taken")
		//return errors.New("email already taken")
	}
	return nil
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)

	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	if !emailRegex.MatchString(user.Email) {
		uv.Errors["email"] = append(uv.Errors["email"], "invalid email")
		//return errors.New("invalid email")
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		uv.Errors["password"] = append(uv.Errors["password"], "password is required")
		//return errors.New("password is required")
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 6 {
		uv.Errors["password"] = append(uv.Errors["password"], "password too short")
		//return errors.New("password too short")
	}
	return nil
}

func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	passwordBytes := []byte(user.Password)
	hashedBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		uv.Errors["password"] = append(uv.Errors["password"], err.Error())
		//return err
		return nil
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) createdAt(user *User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return nil
}

func (uv *userValidator) updatedAt(user *User) error {
	user.UpdatedAt = time.Now()
	return nil
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	u := &User{Email: email}
	if err := runUserValFuncs(
		u,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
	); err != nil {
		return nil, errors2.NewHttpError(err, http.StatusBadRequest, "bad request")
	}
	if len(uv.Errors) > 0 {
		return nil, errors2.NewHttpError2(http.StatusBadRequest, "bad request", uv.Errors)
	}
	return uv.UserDB.ByEmail(u.Email)
}

func (uv *userValidator) Create(u *User) error {
	if err := runUserValFuncs(
		u,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailExists,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.setAvatar,
	); err != nil {
		return errors2.NewHttpError(err, http.StatusBadRequest, "")
	}
	if len(uv.Errors) > 0 {
		return errors2.NewHttpError2(http.StatusBadRequest, "bad request", uv.Errors)
	}
	if err := uv.bcryptPassword(u); err != nil {
		return errors2.NewHttpError(err, http.StatusBadRequest, "bad request")
	}
	return uv.UserDB.Create(u)
}
