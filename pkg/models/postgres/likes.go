package postgres

import (
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models/postgres/user"
	"gorm.io/gorm"
)

type Likes struct {
	gorm.Model
	UserID int
	User   user.User
	PostID int
	Post   Post
}
