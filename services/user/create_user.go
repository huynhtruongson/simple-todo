package user_service

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
	"github.com/huynhtruongson/simple-todo/utils"
	"github.com/huynhtruongson/simple-todo/worker"

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
		opts := []asynq.Option{
			asynq.MaxRetry(10),
			asynq.ProcessIn(time.Second * 10),
			asynq.Queue(worker.CriticalQueue),
		}
		err = s.WorkerClient.DistributeTaskSendVerifyEmail(ctx, &worker.TaskSendVerifyEmailPayload{Username: user.Username}, opts...)
		if err != nil {
			return common.NewInternalError(err, common.InternalErrorMessage, "WorkerClient.DistributeTaskSendVerifyEmail")
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
		return common.NewInvalidRequestError(user_entity.ErrorFullnameIsEmpty, user_entity.ErrorFullnameIsEmpty.Error(), "ValidateUser")
	case user.Username == "":
		return common.NewInvalidRequestError(user_entity.ErrorUsernameIsEmpty, user_entity.ErrorUsernameIsEmpty.Error(), "ValidateUser")
	case len(user.Username) < 6:
		return common.NewInvalidRequestError(user_entity.ErrorInvalidUsernameLength, user_entity.ErrorInvalidUsernameLength.Error(), "ValidateUser")
	case user.Password == "":
		return common.NewInvalidRequestError(user_entity.ErrorPasswordIsEmpty, user_entity.ErrorPasswordIsEmpty.Error(), "ValidateUser")
	case len(user.Password) < 6:
		return common.NewInvalidRequestError(user_entity.ErrorInvalidPasswordLength, user_entity.ErrorInvalidPasswordLength.Error(), "ValidateUser")
	case user.Email == "":
		return common.NewInvalidRequestError(user_entity.ErrorEmailIsEmpty, user_entity.ErrorEmailIsEmpty.Error(), "ValidateUser")
	case !user_entity.IsValidEmail(user.Email):
		return common.NewInvalidRequestError(user_entity.ErrorInvalidEmail, user_entity.ErrorInvalidEmail.Error(), "ValidateUser")
	}
	users, err := s.UserRepo.GetUsersByUsername(ctx, s.DB, user.Username)
	if err != nil {
		return common.NewInternalError(err, common.InternalErrorMessage, "ValidateUser.UserRepo.GetUsersByUsername")
	}
	if len(users) > 0 {
		return common.NewInvalidRequestError(user_entity.ErrorUserNameAlreadyExist, user_entity.ErrorUserNameAlreadyExist.Error(), "ValidateUser")
	}
	users, err = s.UserRepo.GetUsersByEmail(ctx, s.DB, user.Email)
	if err != nil {
		return common.NewInternalError(err, common.InternalErrorMessage, "ValidateUser.UserRepo.GetUsersByEmail")
	}
	if len(users) > 0 {
		return common.NewInvalidRequestError(user_entity.ErrorEmailAlreadyExist, user_entity.ErrorEmailAlreadyExist.Error(), "ValidateUser")
	}
	return nil
}
