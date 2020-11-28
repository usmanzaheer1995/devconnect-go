package postgres

import (
	"errors"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models/postgres"
	"regexp"
	"strings"
	"time"

	"github.com/drexedam/gravatar"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	postgres.Model
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"not null;unique_index" json:"email"`
	Password string `gorm:"not null;-" json:"password"`
	PasswordHash string `gorm:"not null"`
	Avatar   string `json:"avatar"`
}

type UserDB interface {
	ByEmail(email string) (*User, error)
	Find(query models.Query) ([]User, int64, error)

	Create(u *User) []error
	Login(email, password string) (*User, error)
}

type UserService interface {
	UserDB
}

type userService struct {
	db *gorm.DB
}

var _ UserService = &userService{}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db,
	}
}

func countUsers(db *gorm.DB, countC chan int64, errC chan error, doneC chan int) {
	var count int64

	err := db.Model(&User{}).Count(&count).Error
	if err != nil {
		errC <- err
		return
	}
	countC <- count
	doneC <- 1
}

func findUsers(db *gorm.DB, query models.Query, users *[]User, errC chan error, doneC chan int) {
	err := db.Limit(int(query.Limit)).Offset(int(query.Offset)).Find(&users).Error

	if err != nil {
		errC <- err
		return
	}
	doneC <- 1
}

func (us *userService) Find(query models.Query) ([]User, int64, error) {
	errC := make(chan error)
	doneC := make(chan int)
	countC := make(chan int64)

	defer close(errC)
	defer close(doneC)
	defer close(countC)

	var usersList []User
	var count int64

	go countUsers(us.db, countC, errC, doneC)
	go findUsers(us.db, query, &usersList, errC, doneC)

	for n := 2; n > 0; {
		select {
		case err := <-errC:
			return nil, 0, err
		case c := <-countC:
			count = c
		case <-doneC:
			n--
		}
	}

	return usersList, count, nil
}

func (us *userService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

func (us *userService) Create(u *User) []error {
	if err := runUserValFuncs(
		u,
		us.normalizeEmail,
		us.requireEmail,
		us.emailFormat,
		us.emailExists,
		us.passwordRequired,
		us.passwordMinLength,
		us.setAvatar,
	); err != nil {
		return err
	}
	var err error
	if err = us.bcryptPassword(u); err != nil {
		return []error{err}
	}
	err = us.db.Create(u).Error
	if err != nil {
		return []error{err}
	}
	return nil
}

func (us *userService) Login(email, password string) (*User, error) {
	user, err := us.ByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, errors.New("incorrect password")
		default:
			return nil, err
		}
	}

	return user, nil
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) []error {
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

func (us *userService) setAvatar(u *User) error {
	url := gravatar.New(u.Email).
		Size(200).
		Default(gravatar.MysteryMan).
		Rating(gravatar.Pg).
		AvatarURL()
	u.Avatar = url
	return nil
}

func (us *userService) requireEmail(user *User) error {
	if user.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

func (us *userService) emailExists(user *User) error {
	existing, err := us.ByEmail(user.Email)
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if user.ID != existing.ID {
		return errors.New("email already taken")
	}
	return nil
}

func (us *userService) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)

	return nil
}

func (us *userService) emailFormat(user *User) error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email")
	}
	return nil
}

func (us *userService) passwordRequired(user *User) error {
	if user.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (us *userService) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 6 {
		return errors.New("password too short")
	}
	return nil
}

func (us *userService) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	passwordBytes := []byte(user.Password)
	hashedBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (us *userService) createdAt(user *User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return nil
}

func (us *userService) updatedAt(user *User) error {
	user.UpdatedAt = time.Now()
	return nil
}

// first will query using the provided gorm.DB and
// it will get the first item returned and place it
// into dst (which should be a pointer).
// If nothing is found in the query, it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	return db.First(dst).Error
}
