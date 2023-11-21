package user_service

import (
	"context"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
	"github.com/huynhtruongson/simple-todo/utils"

	"github.com/jackc/pgx/v5"
)

func (s *UserService) CreateUser(ctx context.Context, user user_entity.User) (int, error) {
	if err := s.ValidateUser(ctx, user); err != nil {
		return 0, err
	}
	var userID int
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		return userID, common.NewInternalError(err, common.InternalErrorMessage, "UserRepo.CreateUser")
	}
	user.Password = hashedPwd
	if err := lib.ExecTX(ctx, s.DB, func(ctx context.Context, tx pgx.Tx) error {
		id, err := s.UserRepo.CreateUser(ctx, tx, user)
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

func (s *UserService) ValidateUser(ctx context.Context, user user_entity.User) error {
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
	users, err := s.UserRepo.GetUsersByUsername(ctx, s.DB, user.Username)
	if err != nil {
		return common.NewInternalError(err, common.InternalErrorMessage, "ValidateUser.UserRepo.GetUsersByUsername")
	}
	if len(users) > 0 {
		return common.NewInvalidRequestError(nil, user_entity.ErrorUserNameAlreadyExist, "ValidateUser")
	}
	return nil
}
