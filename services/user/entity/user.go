package user_entity

import "github.com/huynhtruongson/simple-todo/common"

type User struct {
	UserID   int    `json:"user_id"`
	FullName string `json:"fullname"`
	Username string `json:"username"`
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
			"password",
		}, []interface{}{
			&u.UserID,
			&u.FullName,
			&u.Username,
			&u.Password,
		}
}

const (
	ErrorUserNameAlreadyExist  = "Username has already existed"
	ErrorFullnameIsEmpty       = "Fullname is empty"
	ErrorUsernameIsEmpty       = "Username is empty"
	ErrorInvalidUsernameLength = "Username length is less than 6 characters"
	ErrorPasswordIsEmpty       = "Password is empty"
	ErrorInvalidPasswordLength = "Password length is less than 6 characters"
)
