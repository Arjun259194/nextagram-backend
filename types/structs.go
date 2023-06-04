package types

import (
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type SuccessResponse struct {
	Status       string      `json:"status"`
	Message      string      `json:"message"`
	ResponseData interface{} `json:"responseData"`
}

func NewErrorResponse(code int, err error, message string) ErrorResponse {
	return ErrorResponse{
		Status:  http.StatusText(code),
		Message: message,
		Error:   err.Error(),
	}
}

func NewSuccessResponse(code int, data interface{}, message string) SuccessResponse {
	return SuccessResponse{
		Status:       http.StatusText(code),
		ResponseData: data,
		Message:      message,
	}
}

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginRequestBody) Validate() error {
	if l.Email == "" {
		return errors.New("email is required")
	}
	if l.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type RegisterRequestBody struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
}

func (r *RegisterRequestBody) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if r.Gender == "" {
		return errors.New("gender is required")
	}
	return nil
}

type UpgradeRouteReqBody struct {
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required"`
	Gender string `json:"gender" validate:"required"`
}
