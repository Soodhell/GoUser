package repositories

import (
	"User/interal/entity"
	"User/pkg/db"
	"errors"
	"fmt"
)

type Repository struct {
	DB *db.Postgres
}

func StartRepository(db *db.Postgres) *Repository {
	return &Repository{db}
}

func (r *Repository) Update(oldEmail string, email string, name string, password string, pathImage string, nameImage string) error {
	if email != "" {
		_, err := r.DB.DB.Exec("update users set email = $1, name = $2, password = $3, path_image = $4, name_image = $5 where email = $6", email, name, password, pathImage, nameImage, oldEmail)
		if err != nil {
			return errors.New(fmt.Sprint(err))
		}
	}

	return nil
}

func (r *Repository) Add(email string, name string, password string, pathImage string, nameImage string) bool {

	_, err := r.DB.DB.Exec("insert into users (email, name, password, path_image, name_image, is_active) values ($1, $2, $3, $4, $5, TRUE)", email, name, password, pathImage, nameImage)
	if err != nil {
		return false
	}
	return true
}

func (r *Repository) Get(email string) (entity.User, error) {
	var getUser entity.User

	err := r.DB.DB.QueryRow("select id, email, name, password, path_image, name_image, roles from users where email = $1", email).Scan(&getUser.ID, &getUser.Email, &getUser.Name, &getUser.Password, &getUser.PathImage, &getUser.NameImage, &getUser.Roles)
	if err != nil {
		return entity.User{}, err
	}

	return getUser, nil
}

func (r *Repository) IsActive(email string, isActive bool) error {

	_, err := r.DB.DB.Exec("update users set is_active = $1 where email = $2", isActive, email)
	return err
}
