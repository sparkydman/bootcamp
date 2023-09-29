package utils

import (
	"bootcamp-api/app/model/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageCode int

const (
	SuccessfulCode MessageCode = iota
	ValidationErrorCode
	BadRequestErrorCode
	UnAuthorizedErrorCode
	ForbiddenErrorCode
	NotFoundErrorCode
	ServerErrorCode
	CredentialsErrorCode
)

var errorMessages = []string{"Successful", "Invalid request data", "Bad request", "User not authorized", "User not allowed to perform this action", "Data not found", "Internal server issue", "Invalid credentials"}

var errorStatus = []string{"SUCCESSFUL", "INVALID_DATA", "BAD_REQUEST", "UNAUTHORIZED", "FORBIDDEN", "DATA_NOT_FOUND", "SERVER_ERROR", "INVALID_CREDENTIAL"}

func (m MessageCode) GetStatus() string {
	return errorStatus[m]
}

func (m MessageCode) GetMessage() string {
	return errorMessages[m]
}

func (m *MessageCode) SetMessage(message string) {
	if *m >= 0 && int(*m) < len(errorMessages) {
		errorMessages[*m] = message
	}
}

func NULL() interface{} {
	return nil
}

func SetResponse[T any](isSuccessful bool, status MessageCode, data T) dto.Response[T] {
	return dto.Response[T]{
		IsSuccessful: isSuccessful,
		Message:      status.GetMessage(),
		Status:       status.GetStatus(),
		Data:         data,
	}
}

func UnAuthorizedResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, SetResponse(false, UnAuthorizedErrorCode, NULL()))
	c.Abort()
}
