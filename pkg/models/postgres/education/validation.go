package education

import (
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/types"
	"net/http"
	"time"
)

type educationValidator struct {
	EducationDB
}

type eduValFunc func(pr *types.EducationRequest) error

func runEduValFuncs(pr *types.EducationRequest, fns ...eduValFunc) error {
	for _, fn := range fns {
		if err := fn(pr); err != nil {
			return err
		}
	}
	return nil
}

func NewEduValidator(edb EducationDB) *educationValidator {
	return &educationValidator{
		EducationDB: edb,
	}
}

func (eduv *educationValidator) educationTimeFormat(edu *types.EducationRequest) error {
	cf, err := time.Parse("2006-01-02", edu.From)
	if err != nil {
		return models.NewHttpError(err, http.StatusBadRequest, "error parsing date")
	}
	edu.ConvertedFrom = cf

	ct, err := time.Parse("2006-01-02", edu.To)
	if err != nil {
		return models.NewHttpError(err, http.StatusBadRequest, "")
	}
	edu.ConvertedTo = ct
	return nil
}

func (eduv *educationValidator) schoolRequired(edu *types.EducationRequest) error {
	if edu.School == "" {
		return &models.HttpError{
			Cause:  nil,
			Detail: "school is required",
			Status: http.StatusBadRequest,
		}
	}
	return nil
}

func (eduv *educationValidator) degreeRequired(edu *types.EducationRequest) error {
	if edu.Degree == "" {
		return &models.HttpError{
			Cause:  nil,
			Detail: "degree is required",
			Status: http.StatusBadRequest,
		}
	}
	return nil
}

func (eduv *educationValidator) fieldofstudyRequired(edu *types.EducationRequest) error {
	if edu.Fieldofstudy == "" {
		return &models.HttpError{
			Cause:  nil,
			Detail: "field of study is required",
			Status: http.StatusBadRequest,
		}
	}
	return nil
}

func (eduv *educationValidator) fromRequired(edu *types.EducationRequest) error {
	if edu.From == "" {
		return &models.HttpError{
			Cause:  nil,
			Detail: "from date is required",
			Status: http.StatusBadRequest,
		}
	}
	return nil
}

func (pv *educationValidator) Create(edu *types.EducationRequest) error {
	if err := runEduValFuncs(
		edu,
		pv.schoolRequired,
		pv.degreeRequired,
		pv.fieldofstudyRequired,
		pv.fromRequired,
		pv.educationTimeFormat,
	); err != nil {
		return models.NewHttpError(err, http.StatusBadRequest, "")
	}
	return pv.EducationDB.Create(edu)
}