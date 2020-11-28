package models

import (
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models/postgres"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID int
	User   postgres.User
	PostID int
	Post   Post
	Text   string `gorm:"not null"`
	Name   string
	Avatar string
}
