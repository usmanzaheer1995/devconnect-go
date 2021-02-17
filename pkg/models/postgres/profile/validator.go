package profile

import (
	"errors"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/types"
	"net/http"
	"strings"
	"time"
)

type profileValidator struct {
	ProfileDB
}

type profileValFunc func(pr *types.ProfileRequest) error

func runProfileValFuncs(pr *types.ProfileRequest, fns ...profileValFunc) error {
	for _, fn := range fns {
		if err := fn(pr); err != nil {
			return err
		}
	}
	return nil
}

func NewProfileValidator(pdb ProfileDB) *profileValidator {
	return &profileValidator{
		ProfileDB: pdb,
	}
}

func (pv *profileValidator) skillsRequired(p *types.ProfileRequest) error {
	if len(p.Skills) == 0 {
		return errors.New("atleast one skill required")
	}
	return nil
}

func (pv *profileValidator) statusRequired(p *types.ProfileRequest) error {
	if p.Status == "" {
		return errors.New("status is required")
	}
	return nil
}

func (pv *profileValidator) skillsFormat(p *types.ProfileRequest) error {
	p.SkillList = strings.Split(p.Skills, ",")
	for i := range p.SkillList {
		p.SkillList[i] = strings.TrimSpace(p.SkillList[i])
	}
	return nil
}

func (pv *profileValidator) experienceTimeFormat(p *types.ProfileRequest) error {
	for idx, _ := range p.Experience {
		cf, err := time.Parse("2006-01-02", p.Experience[idx].From)
		if err != nil {
			return err
		}
		p.Experience[idx].ConvertedFrom = cf

		ct, err := time.Parse("2006-01-02", p.Experience[idx].To)
		if err != nil {
			return err
		}
		p.Experience[idx].ConvertedTo = ct
	}
	return nil
}

func (pv *profileValidator) Create(p *types.ProfileRequest) error {
	if err := runProfileValFuncs(
		p,
		pv.statusRequired,
		pv.skillsRequired,
		pv.skillsFormat,
		pv.experienceTimeFormat,
	); err != nil {
		return models.NewHttpError(err, http.StatusBadRequest, "")
	}
	return pv.ProfileDB.Create(p)
}

func(pv *profileValidator) Update(p *types.ProfileRequest) error {
	if err := runProfileValFuncs(
		p,
		pv.skillsFormat,
		pv.experienceTimeFormat,
	); err != nil {
		return models.NewHttpError(err, http.StatusBadRequest, "")
	}
	return pv.ProfileDB.Update(p)
}
