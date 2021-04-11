package education

import (
	"github.com/usmanzaheer1995/devconnect-go-v2/internal/types"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/utils"
	"gorm.io/gorm"
	"time"
)

type Education struct {
	utils.GormModel
	School       string    `gorm:"not null"`
	Degree       string    `gorm:"not null"`
	Fieldofstudy string    `gorm:"not null"`
	From         time.Time `gorm:"not null;type:time"`
	To           time.Time `gorm:"type:time"`
	Current      bool `gorm:"default:false"`
	Description  string
	ProfileID    uint
}

type EducationDB interface {
	Find(p *Education) (*Education, error)
	FindByProfile(pid uint) ([]Education, error)

	Create(e *types.EducationRequest) error
	Update(e *types.EducationRequest) error
}

type educationService struct {
	EducationDB
}

type EducationService interface {
	EducationDB
}

var _ EducationService = &educationService{}

func NewProfileService(db *gorm.DB) *educationService {
	pg := newEducationGorm(db)
	pv := NewEduValidator(pg)
	return &educationService{
		EducationDB: pv,
	}
}
