package biz

import (
	"github.com/jboolean/whereabouts-api/dao"
	"github.com/jboolean/whereabouts-api/model"
	"log"
)

func FindUserByUsername(username string) (*model.User, error) {
	return dao.UserDAO.FindByUsername(username)
}

func FindRolesByUsername(username string) (model.Roles, error) {
	return dao.UserDAO.FindRoles(username)
}

type DuplicateUsernameError string

func (err DuplicateUsernameError) Error() string {
	return "There is already a user with username: " + string(err)
}

func CreateUser(user *model.User, roles model.Roles, password string) error {
	result, err := dao.UserDAO.FindByUsername(user.Username)
	if err != nil {
		return err
	}
	if result != nil {
		return DuplicateUsernameError(user.Username)
	}

	if password != "" {
		user.Salt, user.Password = createHashedPassword(password)
	}

	log.Printf("Creating user %v", user)

	err = dao.UserDAO.Create(*user)
	if err != nil {
		return err
	}
	err = dao.UserDAO.CreateRoles(user.Username, roles)
	return err
}

func CheckUserPassword(username, password string) bool {
	user, err := FindUserByUsername(username)
	if err != nil {
		log.Panic(err)
	}
	if user == nil {
		return false
	}
	return checkPassword(password, user.Salt, user.Password)
}

func CheckOrSetUserPassword(username, password string) bool {
	user, err := FindUserByUsername(username)
	if err != nil {
		log.Panic(err)
	}
	if user == nil {
		return false
	}
	log.Printf("USER: %#v", user)
	if user.Password == "" {
		user.Salt, user.Password = createHashedPassword(password)
		dao.UserDAO.UpdatePassword(user.Username, user.Password, user.Salt)
	}

	return checkPassword(password, user.Salt, user.Password)
}

func UpdateUserPassword(username, password string) {
	if password == "" {
		dao.UserDAO.UpdatePassword(username, "", "")
		return
	}
	salt, hashedPassword := createHashedPassword(password)
	dao.UserDAO.UpdatePassword(username, hashedPassword, salt)
}

func DeleteUser(username string) {
	dao.UserDAO.Delete(username)
}
