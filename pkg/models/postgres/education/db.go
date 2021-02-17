package education

import (
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/types"
	"gorm.io/gorm"
)

type educationGorm struct {
	db *gorm.DB
}

var _ EducationDB = &educationGorm{}

func newEducationGorm(db *gorm.DB) *educationGorm {
	return &educationGorm{db}
}

func (e educationGorm) Find(p *Education) (*Education, error) {
	panic("implement me")
}

func (e educationGorm) FindByProfile(pid uint) ([]Education, error) {
	panic("implement me")
}

func (e educationGorm) Create(p *types.EducationRequest) error {
	panic("implement me")
}

func (e educationGorm) Update(p *types.EducationRequest) error {
	panic("implement me")
}
