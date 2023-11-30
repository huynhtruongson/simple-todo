package user_entity

import (
	"errors"
	"regexp"

	"github.com/huynhtruongson/simple-todo/common"
)

type User struct {
	UserID   int    `json:"user_id"`
	FullName string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	common.SQLModel
}

func (u User) TableName() string {
	return "users"
}

func (u *User) FieldMap() ([]string, []interface{}) {
	return []string{
			"user_id",
			"fullname",
			"username",
			"email",
			"password",
		}, []interface{}{
			&u.UserID,
			&u.FullName,
			&u.Username,
			&u.Email,
			&u.Password,
		}
}

var (
	ErrorUserNameAlreadyExist  = errors.New("Username has already existed")
	ErrorEmailAlreadyExist     = errors.New("Email has already existed")
	ErrorFullnameIsEmpty       = errors.New("Fullname is empty")
	ErrorUsernameIsEmpty       = errors.New("Username is empty")
	ErrorEmailIsEmpty          = errors.New("Email is empty")
	ErrorInvalidUsernameLength = errors.New("Username length is less than 6 characters")
	ErrorPasswordIsEmpty       = errors.New("Password is empty")
	ErrorInvalidPasswordLength = errors.New("Password length is less than 6 characters")
	ErrorInvalidEmail          = errors.New("Email is invalid")
)

func IsValidEmail(email string) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return regex.MatchString(email)
}
