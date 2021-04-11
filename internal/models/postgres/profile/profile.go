package profile

import (
	"time"

	"github.com/lib/pq"
	"github.com/usmanzaheer1995/devconnect-go-v2/internal/models"
	"github.com/usmanzaheer1995/devconnect-go-v2/internal/models/postgres/user"
	"github.com/usmanzaheer1995/devconnect-go-v2/internal/types"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/utils"
	"gorm.io/gorm"
)

type Experience struct {
	utils.GormModel
	Title       string `gorm:"not null"`
	Company     string `gorm:"not null"`
	Location    string
	From        time.Time
	To          time.Time
	Current     bool `gorm:"default:false"`
	Description string
	ProfileID   uint
}

type Social struct {
	Youtube   string `json:"youtube"`
	Twitter   string `json:"twitter"`
	Facebook  string `json:"facebook"`
	Linkedin  string `json:"linkedin"`
	Instagram string `json:"instagram"`
}

type Profile struct {
	utils.GormModel
	Company        string         `json:"company"`
	Website        string         `json:"website"`
	Location       string         `json:"location"`
	Status         string         `json:"status"`
	Skills         pq.StringArray `gorm:"type:text[]" json:"skills"`
	Bio            string         `json:"bio"`
	Githubusername string         `json:"githubusername"`
	UserID         int
	User           user.User
	Social
}

type ProfileDB interface {
	Find(p *Profile) (*Profile, error)
	FindAll(query models.Query) ([]Profile, int64, error)
	FindByUser(uid uint) (*Profile, error)

	Create(p *types.ProfileRequest) error
	Update(p *types.ProfileRequest) error
}

type profileService struct {
	ProfileDB
}

var _ ProfileService = &profileService{}

type ProfileService interface {
	ProfileDB
}

func NewProfileService(db *gorm.DB) *profileService {
	pg := newProfileGorm(db)
	pv := NewProfileValidator(pg)
	return &profileService{
		ProfileDB: pv,
	}
}
