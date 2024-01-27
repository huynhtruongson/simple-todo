package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	common_resp "github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/features/common"
	"github.com/huynhtruongson/simple-todo/field"
	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
	user_repo "github.com/huynhtruongson/simple-todo/services/user/repository"
	"github.com/huynhtruongson/simple-todo/utils"
	"github.com/stretchr/testify/assert"
)

func (s *suite) createUser(ctx context.Context, payloadType string) (context.Context, error) {
	stepState := common.StepStateFromContext(ctx)
	randStr := utils.RandomString(10)
	userPayload := user_entity.User{
		FullName: field.NewString(fmt.Sprintf("fullname %s", randStr)),
		Username: field.NewString(fmt.Sprintf("username %s", randStr)),
		Email:    field.NewString(fmt.Sprintf("email-%s@gmail.com", randStr)),
		Password: field.NewString(randStr),
	}
	switch payloadType {
	case "empty fullname":
		userPayload.FullName = field.NewNullString()
	case "empty username":
		userPayload.Username = field.NewNullString()
	case "invalid username":
		userPayload.Username = field.NewString("son-1")
	case "existing username":
		userPayload.Username = field.NewString("usernameseed")
	case "empty email":
		userPayload.Email = field.NewNullString()
	case "invalid email":
		userPayload.Email = field.NewString("invalid-email")
	case "existing email":
		userPayload.Email = field.NewString("email+seed@gmail.com")
	case "empty password":
		userPayload.Password = field.NewNullString()
	case "invalid password":
		userPayload.Password = field.NewString("123")
	}

	jsonPayload, err := json.Marshal(userPayload)
	if err != nil {
		return ctx, fmt.Errorf("json marshal user payload failed,%w\n", err)
	}
	statusCode, resBody, err := common.MakeHttpRequest(fmt.Sprintf("%s/%s", s.BaseURL, "v1/user/create"), http.MethodPost, jsonPayload, "")
	if err != nil {
		return ctx, err
	}

	if payloadType == "valid information" {
		err = common.AssertExpectedAndActual(assert.Equal, http.StatusOK, statusCode, "response status code mismatch")
		if err != nil {
			return ctx, err
		}
		var respStruct common_resp.SuccessResponse
		err = json.Unmarshal(resBody, &respStruct)
		if err != nil {
			return ctx, fmt.Errorf("json unmarshal response data failed,%w\n", err)
		}
		userPayload.UserID = field.NewInt(int(respStruct.Data.(float64)))
		stepState.Set(common.RequestPayloadKey, userPayload)
		return common.StepStateToContext(ctx, stepState), nil
	} else {
		var errorStruct common_resp.AppError
		err = json.Unmarshal(resBody, &errorStruct)
		if err != nil {
			return ctx, fmt.Errorf("json unmarshal response data failed,%w\n", err)
		}
		stepState.Set(common.ResponseKey, errorStruct)
		return common.StepStateToContext(ctx, stepState), nil
	}
}

func (s *suite) assertCreatedUser(ctx context.Context) (context.Context, error) {
	stepState := common.StepStateFromContext(ctx)
	payload, ok := stepState.Get(common.RequestPayloadKey).(user_entity.User)
	if !ok {
		return ctx, fmt.Errorf("retrieve user payload failed")
	}
	userRepo := user_repo.NewUserRepo()
	usersDB, err := userRepo.GetUsersByUserIds(ctx, s.DB, []int{payload.UserID.Int()})
	if err != nil {
		return ctx, fmt.Errorf("get user by ids error,%w\n", err)
	}

	if len(usersDB) != 1 {
		return ctx, fmt.Errorf("user not found in db,user_id:%d\n", payload.UserID.Int())
	}
	if err := common.AssertExpectedAndActual(assert.Equal, payload.FullName.String(), usersDB[0].FullName.String()); err != nil {
		return ctx, err
	}
	if err := common.AssertExpectedAndActual(assert.Equal, payload.Email.String(), usersDB[0].Email.String()); err != nil {
		return ctx, err
	}
	return ctx, nil
}

func (s *suite) login(ctx context.Context) (context.Context, error) {
	stepState := common.StepStateFromContext(ctx)
	payload, ok := stepState.Get(common.RequestPayloadKey).(user_entity.User)
	if !ok {
		return ctx, fmt.Errorf("retrieve user payload failed")
	}
	cred := auth_entity.Credential{Username: payload.Username, Password: payload.Password}
	jsonPayload, err := json.Marshal(cred)
	if err != nil {
		return ctx, fmt.Errorf("json marshal user payload failed,%w\n", err)
	}
	statusCode, resBody, err := common.MakeHttpRequest(fmt.Sprintf("%s/%s", s.BaseURL, "v1/auth/login"), http.MethodPost, jsonPayload, "")

	err = common.AssertExpectedAndActual(assert.Equal, http.StatusOK, statusCode, "response status code mismatch")
	if err != nil {
		return ctx, err
	}
	var respStruct auth_entity.LoginResponse
	err = json.Unmarshal(resBody, &respStruct)
	if err != nil {
		return ctx, fmt.Errorf("json unmarshal response data failed,%w\n", err)
	}
	if err := common.AssertActual(assert.NotEmpty, respStruct.AccessToken.String()); err != nil {
		return ctx, err
	}
	if err := common.AssertActual(assert.NotEmpty, respStruct.RefreshToken.String()); err != nil {
		return ctx, err
	}
	return ctx, nil
}
