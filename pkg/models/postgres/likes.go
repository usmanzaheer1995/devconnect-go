package models

import (
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models/postgres"
	"gorm.io/gorm"
)

type Likes struct {
	gorm.Model
	UserID int
	User   postgres.User
	PostID int
	Post   Post
}
