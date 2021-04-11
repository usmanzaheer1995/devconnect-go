package postgres

import (
	"github.com/usmanzaheer1995/devconnect-go-v2/internal/models/postgres/user"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID int
	User   user.User
	PostID int
	Post   Post
	Text   string `gorm:"not null"`
	Name   string
	Avatar string
}
