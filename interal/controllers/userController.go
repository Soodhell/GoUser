package controllers

import (
	"User/interal/DTO"
	"User/interal/entity"
	"User/interal/jwtWork"
	"User/interal/services"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	Service services.Service
}

func StartController(service services.Service) *UserController {
	return &UserController{service}
}

func (u *UserController) SettingRouter(router *mux.Router) {

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/account/sign-up", u.Reg).Methods("POST")
	router.HandleFunc("/account/sign-in", u.Auth).Methods("POST")
	router.HandleFunc("/account/{email}", u.GetUser).Methods("GET")
	router.HandleFunc("/account", u.GetUserJWT).Methods("GET")
	router.HandleFunc("/account/update", u.Update).Methods("PATCH")
	router.HandleFunc("/account/delete", u.UserNotActive).Methods("DELETE")
	router.HandleFunc("/account/recovery", u.UserActive).Methods("POST")
	router.HandleFunc("/account/img/{image}", u.GetImages).Methods("GET")

}

// GetImages godoc
// @Summary Get image
// @Tags Image
// @Description get image
// @ID get-image
// @Produce octet-stream
// @Param image path string true "Image name"
// @Success 200 {string} binary
// @Failure 400 {object} DTO.ResponseError
// @Router /account/img/{image} [GET]
func (u *UserController) GetImages(w http.ResponseWriter, r *http.Request) {

	image := mux.Vars(r)["image"]

	file, err := os.ReadFile("static/avatars/" + image)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(u.errorReturnMessage("такого файла не существует"))
		return
	}

	format := u.detectContentType(file)
	if format == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(u.errorReturnMessage("не приемлемый формат файла"))
		return
	}
	w.Header().Set("Content-Type", format)

	w.Write(file)
	w.WriteHeader(http.StatusOK)

}

// Update godoc
// @Summary Update user
// @Tags User
// @Description update user
// @ID update-user
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param email formData string false "New email"
// @Param password formData string false "New password"
// @Param name formData string false "New name"
// @Param file formData file false "New avatar"
// @Success 200 {object} DTO.ResponseToken
// @Success 200 {object} DTO.ResponseSuccess
// @Failure 400 {object} DTO.ResponseError
// @Failure 401 {object} DTO.ResponseError
// @Router /account/update [PATCH]
func (u *UserController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authToken, err := u.getJWT(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(DTO.ResponseError{err.Error()})
		return
	}

	oldEmail, err := jwtWork.VerifyToken(authToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(DTO.ResponseError{"нет такого пользователя"})
		return
	}

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	name := r.PostFormValue("name")
	file, header, errFile := r.FormFile("file")

	emailUpdate := true
	if email == "" {
		emailUpdate = false
	}

	if errFile != nil && name == "" && password == "" && email == "" {
		json.NewEncoder(w).Encode(DTO.ResponseSuccess{Mail: "ничего не было изменено", Success: true})
		w.WriteHeader(http.StatusOK)
		return
	}

	err = u.Service.Update(oldEmail, email, name, password, file, header, errFile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DTO.ResponseError{err.Error()})
		return
	}

	if emailUpdate {

		token, err := jwtWork.CreateToken(entity.User{Email: email, Password: password, Name: name})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(u.errorReturnMessage(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DTO.ResponseToken{Success: true, Token: token})
		return
	}

	json.NewEncoder(w).Encode(DTO.ResponseSuccess{Success: true, Mail: "Успешно изменено"})
	w.WriteHeader(http.StatusOK)
}

// GetUserJWT godoc
// @Summary Get user by JWT
// @Tags User
// @Description get user info using JWT token
// @ID get-user-jwt
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} DTO.ReturnUser
// @Failure 400 {object} DTO.ResponseError
// @Failure 401 {object} DTO.ResponseError
// @Router /account [GET]
func (u *UserController) GetUserJWT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authToken, err := u.getJWT(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(DTO.ResponseError{err.Error()})
		return
	}

	email, err := jwtWork.VerifyToken(authToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(DTO.ResponseError{"нет такого пользователя"})
		return
	}

	var response DTO.ReturnUser

	user, err := u.Service.GetUserByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DTO.ResponseError{"нет такого пользователя"})
		return
	}

	response.Email = user.Email
	response.Name = user.Name
	response.Id = user.ID
	response.Roles = user.Roles
	response.Image = "account/img/" + user.NameImage

	json.NewEncoder(w).Encode(response)

	w.WriteHeader(http.StatusOK)
}

// Reg godoc
// @Summary Register user
// @Tags User
// @Description register new user
// @ID register-user
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param name formData string true "Name"
// @Param file formData file true "Avatar"
// @Success 200 {object} DTO.ResponseSuccess
// @Failure 400 {object} DTO.ResponseError
// @Router /account/sign-up [POST]
func (u *UserController) Reg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(u.errorReturnMessage("не найден файл"))
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("name")

	err = u.Service.Add(email, name, password, file, header)

	var e *DTO.ErrorUser
	switch {
	case errors.As(err, &e):
		w.WriteHeader(e.Code)
		json.NewEncoder(w).Encode(u.errorReturnMessage(e.Error()))
		return
	}

	defer file.Close()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DTO.ResponseSuccess{Success: true, Mail: "пользователь создан"})
}

// Auth godoc
// @Summary Authenticate user
// @Tags User
// @Description authenticate user and get JWT token
// @ID authenticate-user
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Success 200 {object} DTO.ResponseToken
// @Failure 400 {object} DTO.ResponseError
// @Router /account/sign-in [POST]
func (u *UserController) Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DTO.ResponseError{Error: "не все поля заполнены"})
		return
	}

	var user entity.User
	var err error

	user, err = u.Service.GetUserByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(u.errorReturnMessage("такого пользователя нет"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(u.errorReturnMessage("пароль не верный"))
		return
	} else {
		var resp DTO.ResponseToken
		resp.Success = true
		resp.Token, err = jwtWork.CreateToken(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(u.errorReturnMessage(err.Error()))
			return
		}

		json.NewEncoder(w).Encode(resp)
	}

}

// UserActive godoc
// @Summary Activate user
// @Tags User
// @Description activate user account
// @ID activate-user
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Success 200 {object} DTO.ResponseToken
// @Failure 400 {object} DTO.ResponseError
// @Router /account/recovery [POST]
func (u *UserController) UserActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DTO.ResponseError{Error: "не все поля заполнены"})
	}

	var user entity.User
	var err error

	user, err = u.Service.GetUserByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(u.errorReturnMessage("такого пользователя нет"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(u.errorReturnMessage("пароль не верный"))
		return
	} else {
		var resp DTO.ResponseToken
		resp.Success = true
		resp.Token, err = jwtWork.CreateToken(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(u.errorReturnMessage(err.Error()))
			return
		}

		json.NewEncoder(w).Encode(resp)
		u.Service.IsActive(email, true)
	}

}

// UserNotActive godoc
// @Summary Deactivate user
// @Tags User
// @Description deactivate user account
// @ID deactivate-user
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} DTO.ResponseSuccess
// @Failure 400 {object} DTO.ResponseError
// @Failure 401 {object} DTO.ResponseError
// @Router /account/delete [DELETE]
func (u *UserController) UserNotActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authToken, err := u.getJWT(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(u.errorReturnMessage("не зарегестрирован"))
		return
	}

	email, err := jwtWork.VerifyToken(authToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(u.errorReturnMessage("такого пользователя нет"))
		return
	}

	err = u.Service.IsActive(email, false)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(u.errorReturnMessage("не существует такого пользователя"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DTO.ResponseSuccess{Success: true, Mail: "пользователь успешно удален"})

}

// GetUser godoc
// @Summary Get user by email
// @Tags User
// @Description get user info by email
// @ID get-user-by-email
// @Produce json
// @Param email path string true "User email"
// @Success 200 {object} DTO.ReturnUser
// @Failure 404
// @Router /account/{email} [GET]
func (u *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	email := mux.Vars(r)["email"]

	user, err := u.Service.GetUserByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	returnUser := DTO.ReturnUser{
		Email: user.Email,
		Id:    user.ID,
		Name:  user.Name,
		Image: "account/img/" + user.NameImage,
		Roles: user.Roles,
	}

	err = json.NewEncoder(w).Encode(returnUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (u *UserController) detectContentType(data []byte) string {
	if len(data) > 8 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return "image/png"
	}
	if len(data) > 2 && data[0] == 0xFF && data[1] == 0xD8 {
		return "image/jpeg"
	}
	if len(data) > 6 && string(data[0:6]) == "GIF87a" || string(data[0:6]) == "GIF89a" {
		return "image/gif"
	}
	if len(data) > 12 && string(data[0:12]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return "image/webp"
	}
	return ""
}

func (u *UserController) errorReturnMessage(err string) DTO.ResponseError {

	var e DTO.ResponseError
	e.Error = err

	return e
}

func (u *UserController) getJWT(r *http.Request) (string, error) {
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		return "", errors.New("не зарегестрирован")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authToken, prefix) {
		return "", errors.New("нет такого пользователя")
	}

	authToken = strings.TrimPrefix(authToken, prefix)
	return authToken, nil
}
