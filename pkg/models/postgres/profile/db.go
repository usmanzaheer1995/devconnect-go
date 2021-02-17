package profile

import (
	"errors"
	"fmt"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/types"
	"gorm.io/gorm"
	"net/http"
)

type profileGorm struct {
	db *gorm.DB
}

var _ ProfileDB = &profileGorm{}

func newProfileGorm(db *gorm.DB) *profileGorm {
	return &profileGorm{db}
}

func (pg *profileGorm) Find(p *Profile) (*Profile, error) {
	return nil, nil
}

func (pg *profileGorm) FindAll(query models.Query) ([]Profile, int64, error) {
	return nil, 0, nil
}

func (pg *profileGorm) FindByUser(uid uint) (*Profile, error) {
	var p Profile
	db := pg.db.Where("user_id = ?", uid)
	err := db.First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (pg *profileGorm) Create(p *types.ProfileRequest) error {
	existingProfile, err := pg.FindByUser(uint(p.UserID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error updating profile: %v", err)
	}
	if existingProfile != nil {
		return models.NewHttpError(nil, http.StatusConflict, "profile already exists")
	}
	tx := pg.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	newProfile := Profile{
		Company:        p.Company,
		Website:        p.Website,
		Location:       p.Location,
		Status:         p.Status,
		Skills:         p.SkillList,
		Bio:            p.Bio,
		Githubusername: p.Githubusername,
		UserID:         p.UserID,
	}

	if err := pg.db.Create(&newProfile).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error creating new profile: %v", err)
	}

	for _, exp := range p.Experience {
		newExp := Experience{
			Title:       exp.Title,
			Company:     exp.Company,
			Location:    exp.Location,
			From:        exp.ConvertedFrom,
			To:          exp.ConvertedTo,
			Current:     exp.Current,
			Description: exp.Description,
			ProfileID:   newProfile.ID,
		}
		if err := pg.db.Create(&newExp).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error creating new experience: %v", err)
		}
	}

	social := Social{
		Youtube:   p.Social.Youtube,
		Twitter:   p.Social.Twitter,
		Facebook:  p.Social.Facebook,
		Linkedin:  p.Social.Linkedin,
		Instagram: p.Social.Instagram,
		ProfileID: newProfile.ID,
	}
	if err := pg.db.Create(&social).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error creating new social: %v", err)
	}

	return nil
}

func (pg *profileGorm) Update(p *types.ProfileRequest) error {
	profile, err := pg.FindByUser(uint(p.UserID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NewHttpError(nil, http.StatusNotFound, "profile not found")
		}
		return fmt.Errorf("error updating profile: %v", err)
	}
	tx := pg.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	updates := &Profile{
		Company:        p.Company,
		Website:        p.Website,
		Location:       p.Location,
		Status:         p.Status,
		Skills:         p.SkillList,
		Bio:            p.Bio,
		Githubusername: p.Githubusername,
	}
	err = pg.db.Model(&profile).Where("id = ?", profile.ID).Updates(updates).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating profile: %v", err)
	}

	return nil
}

func updateExp(db *gorm.DB, profileID uint, updates []types.Experience) error {
	var exp Experience
	err := db.Where("profile_id = ?", profileID).Find(&exp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	fmt.Printf("%+v", exp)
	return nil
}

//func updateEdu(db *gorm.DB, profileID uint, updates []types.Education) error {
//	var edu []Education
//	err := db.Where("profile_id = ?", profileID).Find(&edu).Error
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			return nil
//		}
//		return err
//	}
//	fmt.Printf("Length: %d\n", len(edu))
//	return nil
//}
