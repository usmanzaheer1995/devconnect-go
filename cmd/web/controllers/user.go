package controllers

import (
	"encoding/json"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models"
	userModel "github.com/usmanzaheer1995/devconnect-go-v2/pkg/models/postgres/user"
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
func (u *UserController) Find(w http.ResponseWriter, r *http.Request) {
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
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp := map[string]interface{}{
		"users": users,
		"total": count,
	}

	utils.JSON(w, http.StatusOK, &utils.Response{
		StatusCode: http.StatusOK,
		Message:    "users fetched successfully",
		Data:       resp,
	})
}

func (u *UserController) Create(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := &userModel.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	errs := u.us.Create(user)
	if errs != nil {
		var e []string
		for _, err := range errs {
			e = append(e, err.Error())
		}
		utils.JSON(w, http.StatusBadRequest, &utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "bad request",
			Data:       e,
		})
		return
	}

	token, err := utils.EncodeAuthToken(user.ID)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	user.PasswordHash = ""
	utils.JSON(w, http.StatusOK, &utils.Response{
		StatusCode: http.StatusCreated,
		Message:    "user created successfully",
		Data:       map[string]interface{}{
			"user":  user,
			"token": token,
		},
	})
	return
}

func (u *UserController) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := &userModel.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user, err = u.us.Login(user.Email, user.Password)
	if err != nil {
		utils.ERROR(w, http.StatusNotFound, err)
		return
	}
	token, err := utils.EncodeAuthToken(user.ID)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	user.PasswordHash = ""
	utils.JSON(w, http.StatusOK, &utils.Response{
		StatusCode: http.StatusCreated,
		Message:    "user created successfully",
		Data:       map[string]interface{}{
			"user": user,
			"token": token,
		},
	})
	return
}
