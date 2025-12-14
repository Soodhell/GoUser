package services

import (
	"User/interal/DTO"
	"User/interal/entity"
	"User/interal/repositories"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Service struct {
	Repository repositories.Repository
}

func StartService(repository repositories.Repository) *Service {
	return &Service{repository}
}

func (s *Service) saveFile(file multipart.File, header *multipart.FileHeader) (DTO.FileReturn, *DTO.ErrorUser) {
	avatarName := uuid.New().String()
	ext := filepath.Ext(header.Filename)

	if ext == "" {
		ext = ".jpg"
	}

	avatarPath := fmt.Sprintf("static/avatars/%s%s", avatarName, ext)

	dst, err := os.Create(avatarPath)
	if err != nil {
		return DTO.FileReturn{}, &DTO.ErrorUser{
			Code: http.StatusInternalServerError,
			Msg:  "что то не так на стороне сервера",
		}
	}
	defer dst.Close()
	defer file.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return DTO.FileReturn{}, &DTO.ErrorUser{
			Code: http.StatusInternalServerError,
			Msg:  "что то не так на стороне сервера",
		}
	}

	return DTO.FileReturn{AvatarName: avatarName + ext, AvatarPath: avatarPath}, nil
}

func (s *Service) validFile(header *multipart.FileHeader) *DTO.ErrorUser {

	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		return &DTO.ErrorUser{
			Code: http.StatusBadRequest,
			Msg:  "не верный формат файла (разрешено только png, svg, webp, jpg)",
		}
	}

	err := os.MkdirAll("static/avatars", 0755)
	if err != nil {
		return &DTO.ErrorUser{
			Code: http.StatusInternalServerError,
			Msg:  "что то не так на стороне сервера",
		}
	}

	return nil

}

func (s *Service) Add(email string, name string, password string, file multipart.File, header *multipart.FileHeader) error {

	if email == "" || name == "" || password == "" {
		return &DTO.ErrorUser{
			Code: http.StatusBadRequest,
			Msg:  "не все поля заполнены",
		}
	}

	err := s.validFile(header)
	if err != nil {
		return err
	}

	fileData, errSaveFile := s.saveFile(file, header)
	if errSaveFile != nil {
		return errSaveFile
	}

	bcryptPassword, errGeneratePassword := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errGeneratePassword != nil {
		return &DTO.ErrorUser{
			Code: http.StatusInternalServerError,
			Msg:  "что то не так со стороны сервера",
		}
	}

	if !s.Repository.Add(email, name, string(bcryptPassword), fileData.AvatarPath, fileData.AvatarName) {
		errFileRemove := os.Remove(fileData.AvatarPath)
		if errFileRemove != nil {
			panic("Что то не так с удалением файла - " + errFileRemove.Error() + ". Название файла: " + fileData.AvatarPath)
		}
		return &DTO.ErrorUser{
			Code: http.StatusForbidden,
			Msg:  "такой пользователь уже существует",
		}
	}
	return nil
}

func (s *Service) Update(oldEmail string, email string, name string, password string, file multipart.File, header *multipart.FileHeader, fileError error) error {

	user, err := s.GetUserByEmail(oldEmail)
	pathImage := user.PathImage
	nameImage := user.NameImage
	if err != nil {
		return err
	}

	if email == "" {
		email = user.Email
	}
	if name == "" {
		name = user.Name
	}
	if password == "" {
		password = user.Password
	}

	if fileError == nil {
		errFileValid := s.validFile(header)
		if errFileValid != nil {
			return errFileValid
		}

		fileData, errFileSave := s.saveFile(file, header)
		if errFileSave != nil {
			return err
		}

		errFileRemove := os.Remove(pathImage)
		if errFileRemove != nil {
			panic("что то не так с удалением файла - " + errFileRemove.Error() + ". Название файла: " + pathImage)
		}

		pathImage = fileData.AvatarPath
		nameImage = fileData.AvatarName
	}

	err = s.Repository.Update(oldEmail, email, name, password, pathImage, nameImage)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) IsActive(email string, isActive bool) error {
	_, err := s.GetUserByEmail(email)
	if err != nil {
		return err
	}
	return s.Repository.IsActive(email, isActive)
}

func (s *Service) GetUserByEmail(email string) (entity.User, error) {
	return s.Repository.Get(email)
}
