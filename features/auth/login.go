package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	common_resp "github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/features/common"
	"github.com/huynhtruongson/simple-todo/field"
	"github.com/huynhtruongson/simple-todo/pb"
	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"
	"github.com/huynhtruongson/simple-todo/token"
	"github.com/huynhtruongson/simple-todo/utils"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	userIdKey = "userIdKey"
)

func (s *suite) login(ctx context.Context, payloadType, methodType string) (context.Context, error) {
	stepState := common.StepStateFromContext(ctx)
	userReq, ok := stepState.Get(common.RequestPayloadKey).(*pb.CreateUserRequest)
	if !ok {
		return ctx, fmt.Errorf("cannot retrieve user request")
	}
	username, password := userReq.Username, userReq.Password
	switch payloadType {
	case "empty username":
		username = ""
	case "empty password":
		password = ""
	case "incorrect username":
		username = "incorrect-username"
	case "incorrect password":
		password = "incorrect-password"

	}
	if methodType == "grpc" {
		grpcClient := pb.NewAuthServiceClient(s.GrpcConn)
		resp, err := grpcClient.Login(ctx, &pb.LoginRequest{Username: username, Password: password})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.InvalidArgument:
					stepState.Set(common.ResponseKey, common_resp.AppError{Code: http.StatusBadRequest, Message: e.Message()})
				default:
					return ctx, fmt.Errorf("expect grpc status code is %d or %d, but actual is %d\n", codes.OK, codes.InvalidArgument, e.Code())
				}
			} else {
				return ctx, fmt.Errorf("cannot call login by grpc:%w\n", err)
			}
		} else {
			stepState.Set(common.ResponseKey, resp)
		}
	} else {
		cred := auth_entity.Credential{Username: field.NewString(username), Password: field.NewString(password)}
		jsonPayload, err := json.Marshal(cred)
		if err != nil {
			return ctx, fmt.Errorf("json marshal user payload failed,%w\n", err)
		}
		statusCode, resBody, err := common.MakeHttpRequest(fmt.Sprintf("%s/%s", s.BaseURL, "v1/auth/login"), http.MethodPost, jsonPayload, "")
		if statusCode == http.StatusOK {
			var respStruct auth_entity.LoginResponse
			err = json.Unmarshal(resBody, &respStruct)
			if err != nil {
				return ctx, fmt.Errorf("json unmarshal response data failed,%w\n", err)
			}
			stepState.Set(common.ResponseKey, respStruct)
		} else {
			var respStruct common_resp.AppError
			err = json.Unmarshal(resBody, &respStruct)
			if err != nil {
				return ctx, fmt.Errorf("json unmarshal response data failed,%w\n", err)
			}
			stepState.Set(common.ResponseKey, respStruct)
		}
	}
	return common.StepStateToContext(ctx, stepState), nil
}

func (s *suite) assertToken(ctx context.Context, methodType string) (context.Context, error) {
	stepState := common.StepStateFromContext(ctx)
	acToken, rfToken := "", ""
	if methodType == "grpc" {
		resp, ok := stepState.Get(common.ResponseKey).(*pb.LoginResponse)
		if !ok {
			return ctx, fmt.Errorf("retrieve login response failed")
		}
		if !ok {
			return ctx, fmt.Errorf("retrieve login response failed")
		}
		acToken = resp.AccessToken
		rfToken = resp.RefreshToken
	} else {
		resp, ok := stepState.Get(common.ResponseKey).(auth_entity.LoginResponse)
		if !ok {
			return ctx, fmt.Errorf("retrieve login response failed")
		}
		acToken = resp.AccessToken.String()
		rfToken = resp.RefreshToken.String()
	}
	tokenMaker, err := token.NewPasetoMaker(s.Config.TokenKey)
	if err != nil {
		return ctx, fmt.Errorf("init token maker failed:%w\n", err)
	}
	userId := stepState.Get(userIdKey).(int)
	acPayload, err := tokenMaker.VerifyToken(acToken)
	if err != nil {
		return ctx, fmt.Errorf("verify access token failed:%w\n", err)
	}

	if err := common.AssertExpectedAndActual(assert.Equal, userId, acPayload.UserID); err != nil {
		return ctx, fmt.Errorf("user_id of token payload and created user mismatch :%w\n", err)
	}
	if err := common.AssertWithinDuration(assert.WithinDuration, time.Now(), acPayload.ExpiresAt, time.Minute*15); err != nil {
		return ctx, fmt.Errorf("token expiration time is not within 15 mins :%w\n", err)
	}
	rfPayload, err := tokenMaker.VerifyToken(rfToken)
	if err != nil {
		return ctx, fmt.Errorf("verify refresh token failed:%w\n", err)
	}
	if err := common.AssertExpectedAndActual(assert.Equal, userId, rfPayload.UserID); err != nil {
		return ctx, fmt.Errorf("user_id of token payload and created user mismatch :%w\n", err)
	}
	if err := common.AssertWithinDuration(assert.WithinDuration, time.Now(), rfPayload.ExpiresAt, time.Hour*24); err != nil {
		return ctx, fmt.Errorf("token expiration time is not within 15 mins :%w\n", err)
	}
	return ctx, nil
}

func (s *suite) createUser(ctx context.Context) (context.Context, error) {
	stepState := common.StepStateFromContext(ctx)
	randStr := utils.RandomString(10)
	userPayload := &pb.CreateUserRequest{
		Fullname: fmt.Sprintf("full %s", randStr),
		Email:    fmt.Sprintf("email-%s@gmail.com", randStr),
		Username: fmt.Sprintf("username %s", randStr),
		Password: randStr,
	}
	userSvClient := pb.NewUserServiceClient(s.GrpcConn)
	resp, err := userSvClient.CreateUser(ctx, userPayload)
	if err != nil {
		return ctx, fmt.Errorf("create user by GRPC failed:%w\n", err)
	}
	stepState.Set(common.RequestPayloadKey, userPayload)
	stepState.Set(userIdKey, int(resp.Data))

	return common.StepStateToContext(ctx, stepState), nil
}
