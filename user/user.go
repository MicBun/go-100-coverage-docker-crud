package user

import (
	"crypto/md5"
	"fmt"
	"github.com/MicBun/go-microservice-kubernetes/core"

	"gorm.io/gorm"
)

type Auth struct {
	db *gorm.DB
}

type AuthInterface interface {
	RegisterUser(username, password, name string) (core.User, error)
	AuthenticateUser(username, password string) (core.User, error)
	GetUser(id uint) (core.User, error)
	UpdateUser(id uint, username, password, name string) (core.User, error)
	DeleteUser(id uint) error
	ListUsers() ([]core.User, error)
	SaveToken(id uint, token string) error
}

func AdminAuth(db *gorm.DB) AuthInterface {
	return &Auth{
		db: db,
	}
}

func hash(plain string) string {
	data := []byte(plain)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func (a *Auth) RegisterUser(username, password, name string) (core.User, error) {
	user := core.User{
		Username: username,
		Password: hash(password),
		Name:     name,
	}
	err := a.db.Save(&user).Error
	return user, err
}

func (a *Auth) AuthenticateUser(username, password string) (core.User, error) {
	var user core.User
	err := a.db.Model(core.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return user, fmt.Errorf("unable to retrieve user %w", err)
	}

	h := hash(password)
	if h != user.Password {
		return user, fmt.Errorf("password mismatch")
	}

	return user, nil
}

func (a *Auth) GetUser(id uint) (core.User, error) {
	var user core.User
	err := a.db.Model(core.User{}).Where("id = ?", id).First(&user).Error
	return user, err
}

func (a *Auth) UpdateUser(id uint, username, password, name string) (core.User, error) {
	user, err := a.GetUser(id)
	if err != nil {
		return user, err
	}

	if username != "" {
		user.Username = username
	}
	if password != "" {
		user.Password = hash(password)
	}
	if name != "" {
		user.Name = name
	}
	err = a.db.Updates(&user).Error
	return user, err
}

func (a *Auth) DeleteUser(id uint) error {
	user, err := a.GetUser(id)
	if err != nil {
		return err
	}
	return a.db.Delete(&user).Error
}

func (a *Auth) ListUsers() ([]core.User, error) {
	var users []core.User
	err := a.db.Model(core.User{}).Find(&users).Error
	if len(users) == 0 {
		return users, fmt.Errorf("no users found")
	}
	return users, err
}

func (a *Auth) SaveToken(id uint, token string) error {
	user, err := a.GetUser(id)
	if err != nil {
		return err
	}
	user.Token = token
	return a.db.Save(&user).Error
}
