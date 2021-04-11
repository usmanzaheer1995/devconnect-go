package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	errors2 "github.com/usmanzaheer1995/devconnect-go-v2/internal/errors"
	"github.com/usmanzaheer1995/devconnect-go-v2/internal/models"
	userModel "github.com/usmanzaheer1995/devconnect-go-v2/internal/models/postgres/user"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/utils"
)

type UserController struct {
	us userModel.UserService
}

func NewUserController(u userModel.UserService) *UserController {
	return &UserController{us: u}
}

func (u *UserController) FindByID(w http.ResponseWriter, r *http.Request) error {
	uid := uint(r.Context().Value("userID").(float64))
	user, err := u.us.ByID(uid)
	if err != nil {
		return errors2.NewHttpError(err, http.StatusInternalServerError, "user not found")
	}
	user.PasswordHash = ""
	utils.JSON(w, http.StatusOK, &utils.Response{
		Message: "user fetched successfully",
		Data:    user,
	})
	return nil
}

func (u *UserController) Find(w http.ResponseWriter, r *http.Request) error {
	qp := r.URL.Query()
	var query models.Query
	if qp.Get("offset") != "" {
		off, err := strconv.Atoi(qp.Get("offset"))
		if err != nil {
			query.Offset = 0
		}
		query.Offset = int64(off)
	}
	if qp.Get("limit") != "" {
		lim, err := strconv.Atoi(qp.Get("limit"))
		if err != nil {
			query.Limit = 10
		}
		query.Limit = int64(lim)
	}
	users, count, err := u.us.Find(query)
	if err != nil {
		return fmt.Errorf("error finding users: %v", err)
	}

	utils.JSON(w, http.StatusOK, &utils.Response{
		Message: "users fetched successfully",
		Data: map[string]interface{}{
			"users": users,
			"total": count,
		},
	})
	return nil
}

func (u *UserController) Create(w http.ResponseWriter, r *http.Request) error {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//utils.ERROR(w, http.StatusBadRequest, err)
		return errors.New("Error not null")
	}

	user := &userModel.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return errors.New("error not null")
	}

	err = u.us.Create(user)
	if err != nil {
		return err
	}

	token, err := utils.EncodeAuthToken(user.ID)
	if err != nil {
		return fmt.Errorf("error encoding token: %v", err)
	}

	user.PasswordHash = ""
	utils.JSON(w, http.StatusOK, &utils.Response{
		Message: "user created successfully",
		Data: map[string]interface{}{
			"user":  user,
			"token": token,
		},
	})
	return nil
}

func (u *UserController) Login(w http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	user := &userModel.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return errors2.NewHttpError(err, http.StatusBadRequest, "bad request")
	}

	user, err = u.us.Login(user.Email, user.Password)
	if err != nil {
		return err
	}
	token, err := utils.EncodeAuthToken(user.ID)
	if err != nil {
		return fmt.Errorf("error encoding token: %v", err)
	}

	user.PasswordHash = ""
	utils.JSON(w, http.StatusFound, &utils.Response{
		Message: "user logged in successfully!",
		Data: map[string]interface{}{
			"user":  user,
			"token": token,
		},
	})
	return nil
}
