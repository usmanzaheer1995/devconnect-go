package user

import (
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models"
	"gorm.io/gorm"
)

/* userGorm */
type userGorm struct {
	db *gorm.DB
}

var _ UserDB = &userGorm{}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) Create(u *User) []error {
	err := ug.db.Create(u).Error
	if err != nil {
		return []error{err}
	}
	return nil
}

func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
	return ug.db.Delete(&User{}, id).Error
}

func (us *userGorm) Find(query models.Query) ([]User, int64, error) {
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
/* userGorm */

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

// first will query using the provided gorm.DB and
// it will get the first item returned and place it
// into dst (which should be a pointer).
// If nothing is found in the query, it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	return db.First(dst).Error
}
