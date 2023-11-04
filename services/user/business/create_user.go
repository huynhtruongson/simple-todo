package user_biz

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sondev/todo-list/common"
	"github.com/sondev/todo-list/lib"
	user_entity "github.com/sondev/todo-list/services/user/entity"
)

type UserRepo interface {
	CreateUser(ctx context.Context, db lib.QueryExecer, user user_entity.User) (int, error)
	GetUsersByUsername(ctx context.Context, db lib.QueryExecer, username string) ([]user_entity.User, error)
}
type CreateUserBiz struct {
	DB lib.DB
	UserRepo
}

func NewCreateUserBiz(db lib.DB, userRepo UserRepo) *CreateUserBiz {
	return &CreateUserBiz{
		DB:       db,
		UserRepo: userRepo,
	}
}

func (biz CreateUserBiz) CreateUser(ctx context.Context, user user_entity.User) (int, error) {
	if err := biz.ValidateUser(ctx, user); err != nil {
		return 0, err
	}
	var userID int
	if err := lib.ExecTX(ctx, biz.DB, func(ctx context.Context, tx pgx.Tx) error {
		id, err := biz.UserRepo.CreateUser(ctx, tx, user)
		userID = id
		if err != nil {
			return common.NewInternalError(err, common.InternalErrorMessage, "UserRepo.CreateUser")
		}
		return nil
	}); err != nil {
		return userID, err
	}

	return userID, nil
}

func (biz CreateUserBiz) ValidateUser(ctx context.Context, user user_entity.User) error {
	switch {
	case user.FullName == "":
		return common.NewInvalidRequestError(nil, user_entity.ErrorFullnameIsEmpty, "ValidateUser")
	case user.Username == "":
		return common.NewInvalidRequestError(nil, user_entity.ErrorUsernameIsEmpty, "ValidateUser")
	case len(user.Username) < 6:
		return common.NewInvalidRequestError(nil, user_entity.ErrorInvalidUsernameLength, "ValidateUser")
	case user.Password == "":
		return common.NewInvalidRequestError(nil, user_entity.ErrorPasswordIsEmpty, "ValidateUser")
	case len(user.Password) < 6:
		return common.NewInvalidRequestError(nil, user_entity.ErrorInvalidPasswordLength, "ValidateUser")
	}
	users, err := biz.UserRepo.GetUsersByUsername(ctx, biz.DB, user.Username)
	if err != nil {
		return common.NewInternalError(err, common.InternalErrorMessage, "ValidateUser.UserRepo.GetUsersByUsername")
	}
	if len(users) > 0 {
		return common.NewInvalidRequestError(nil, user_entity.ErrorUserNameAlreadyExist, "ValidateUser")
	}
	return nil
}
