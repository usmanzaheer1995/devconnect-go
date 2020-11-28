package models

import (
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models/postgres"
	"gorm.io/gorm"
	"time"
)

type experience struct {
	Title string `gorm:"not null"`
	Company string `gorm:"not null"`
	Location string
	From time.Time
	To time.Time
	Current bool `gorm:"default:false"`
	Description string
}

type education struct {
	School string `gorm:"not null"`
	Degree string `gorm:"not null"`
	Fieldofstudy string `gorm:"not null"`
	From time.Time `gorm:"not null"`
	To time.Time
	Current bool `gorm:"default:false"`
	Description string
}

type social struct {
	Youtube string
	Twitter string
	Facebook string
	Linkedin string
	Instagram string
}

type Profile struct {
	gorm.Model
	Company string
	Website        string
	Location       string
	Status         string
	Skills         []string
	Bio            string
	Githubusername string
	Experience     experience
	Education      education
	Social         social
	UserID         int
	User           postgres.User
}