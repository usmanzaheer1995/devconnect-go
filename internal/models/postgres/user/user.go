package user

import (
	"fmt"
	"github.com/usmanzaheer1995/devconnect-go-v2/internal/errors"
	"github.com/usmanzaheer1995/devconnect-go-v2/internal/models"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type User struct {
	utils.GormModel
	Name         string `gorm:"not null" json:"name"`
	Email        string `gorm:"not null;unique_index" json:"email"`
	Password     string `gorm:"not null;-" json:"password"`
	PasswordHash string `gorm:"not null" json:"-"`
	Avatar       string `json:"avatar"`
}

type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	Find(query models.Query) ([]User, int64, error)

	Create(u *User) error
}

type UserService interface {
	UserDB
	Login(email, password string) (*User, error)
}

type userService struct {
	UserDB
}

var _ UserService = &userService{}

func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{db}
	uv := newUserValidator(ug)
	return &userService{
		UserDB: uv,
	}
}

func (us *userService) Login(email, password string) (*User, error) {
	user, err := us.ByEmail(email)
	if err != nil {
		return nil, errors.NewHttpError(err, http.StatusBadRequest, "incorrect username/password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, errors.NewHttpError(nil, http.StatusBadRequest, "incorrect username/password")
		default:
			return nil, fmt.Errorf("error comparing hash: %v", err)
		}
	}

	return user, nil
}
