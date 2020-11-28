package postgres

import (
	"fmt"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models/postgres/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServicesConfig func(*Services) error

func WithGorm(connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{})
		if err != nil {
			return err
		}
		s.db = db
		fmt.Println("Successfully connected to database!")
		return nil
	}
}

func WithUser() ServicesConfig {
	return func(s *Services) error {
		s.User = user.NewUserService(s.db)
		return nil
	}
}

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

// Close closes database connection
func (s *Services) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

// TODO: Add this back when appropriate
//DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.Migrator().DropTable(&user.User{})
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// TODO: Add this back when appropriate
// AutoMigrate with attempt to automatically
// migrate all
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&user.User{})
}

type Services struct {
	db   *gorm.DB
	User user.UserService
}
