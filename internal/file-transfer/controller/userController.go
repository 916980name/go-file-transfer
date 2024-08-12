package controller

import (
	"context"
	"encoding/json"
	"file-transfer/internal/file-transfer/service"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/common"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/model"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return UserController{
		service: service,
	}
}

func (uc *UserController) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Create an instance of the custom type
	var loginRequest v1.UserLoginRequest

	// Unmarshal the request body into the custom type
	err = json.Unmarshal(body, &loginRequest)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	user, err := uc.service.Login(ctx, loginRequest)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}
	loginSucResponse(ctx, w, user)
}

func (uc *UserController) LoginShare(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userId := ctx.Value(common.Trace_request_uid{}).(string)
	url, err := uc.service.CreateLoginUrl(ctx, userId)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}
	errno.WriteResponse(ctx, w, url)
}

func (uc *UserController) LoginByShareLink(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["loginKey"]
	if len(key) < 1 {
		errno.WriteErrorResponse(ctx, w, errno.ErrInvalidParameter)
		return
	}

	user, err := uc.service.LoginByLoginUrl(ctx, key)
	if err != nil {
		errno.WriteErrorResponse(ctx, w, err)
		return
	}
	loginSucResponse(ctx, w, user)
}

func (uc *UserController) UserMe(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	errno.WriteResponse(ctx, w, nil)
}

func loginSucResponse(ctx context.Context, w http.ResponseWriter, user *model.UserInfo) {
	resp := &v1.UserLoginResponse{
		Username:   user.Username,
		Privileges: "P_USER",
		IdKey:      user.Id,
	}
	errno.WriteResponse(ctx, w, resp)
}
