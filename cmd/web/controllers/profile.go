package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/usmanzaheer1995/devconnect-go-v2/internal/models/postgres/user"
	"github.com/usmanzaheer1995/devconnect-go-v2/internal/types"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/utils"
	"io/ioutil"
	"net/http"
)
import "github.com/usmanzaheer1995/devconnect-go-v2/internal/models/postgres/profile"

type ProfileController struct {
	us user.UserService
	ps profile.ProfileService
}

func NewProfileController(u user.UserService, p profile.ProfileService) *ProfileController {
	return &ProfileController{us: u, ps: p}
}

func (pc *ProfileController) Create(w http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error parsing body: %v", err)
	}
	p := &types.ProfileRequest{}
	p.UserID = int(r.Context().Value("userID").(float64))
	if err = json.Unmarshal(body, &p); err != nil {
		return fmt.Errorf("error parsing body: %v", err)
	}

	err = pc.ps.Create(p)
	if err != nil {
		return err
	}
	utils.JSON(w, http.StatusCreated, &utils.Response{
		Message: "profile created successfully",
		Data:    p,
	})
	return nil
}

func (pc *ProfileController) Update(w http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error parsing body: %v", err)
	}
	p := &types.ProfileRequest{}
	p.UserID = int(r.Context().Value("userID").(float64))
	if err = json.Unmarshal(body, &p); err != nil {
		return fmt.Errorf("error parsing body: %v", err)
	}

	err = pc.ps.Update(p)
	if err != nil {
		return fmt.Errorf("error during updating: %v", err)
	}
	utils.JSON(w, http.StatusOK, &utils.Response{
		Message: "profile updated successfully",
		Data:    nil,
	})
	return nil
}
