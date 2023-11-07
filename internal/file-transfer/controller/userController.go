package controller

import (
	"context"
	"encoding/json"
	"file-transfer/internal/file-transfer/service"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/token"
	"io"
	"net/http"
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
		errno.WriteResponse(ctx, w, err, nil)
		return
	}
	tokenStr, err := token.Sign(user.Id, user.Username)
	if err != nil {
		errno.WriteResponse(ctx, w, err, nil)
		return
	}
	w.Header().Set(token.AuthHeader, tokenStr)
	errno.WriteResponse(ctx, w, nil, user.Username)
}
