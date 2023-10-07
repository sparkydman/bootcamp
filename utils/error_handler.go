package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func PanicException(code MessageCode) {
	err := errors.New(code.GetMessage())
	err = fmt.Errorf("%s:%w", code.GetStatus(), err)
	if err != nil {
		panic(err)
	}
}

func ValidationException(err error, code MessageCode) {
	err = fmt.Errorf("%s:%w", code.GetStatus(), err)
	if err != nil {
		panic(err)
	}
}

func ResponseErrorHandler(c *gin.Context) {
	if err := recover(); err != nil {
		errorToString := fmt.Sprint(err)
		strArr := strings.Split(errorToString, ":")
		key := strArr[0]

		switch key {
		case NotFoundErrorCode.GetStatus():
			c.JSON(http.StatusNotFound, SetResponse(false, NotFoundErrorCode, NULL()))
			c.Abort()
		case ValidationErrorCode.GetStatus():
			c.JSON(http.StatusBadRequest, SetResponse(false, ValidationErrorCode, NULL()))
			c.Abort()
		case BadRequestErrorCode.GetStatus():
			c.JSON(http.StatusBadRequest, SetResponse(false, BadRequestErrorCode, NULL()))
			c.Abort()
		case UnAuthorizedErrorCode.GetStatus():
			c.JSON(http.StatusUnauthorized, SetResponse(false, UnAuthorizedErrorCode, NULL()))
			c.Abort()
		case ForbiddenErrorCode.GetStatus():
			c.JSON(http.StatusForbidden, SetResponse(false, ForbiddenErrorCode, NULL()))
			c.Abort()
		case CredentialsErrorCode.GetStatus():
			c.JSON(http.StatusBadRequest, SetResponse(false, CredentialsErrorCode, NULL()))
			c.Abort()
		case ConflictErrorCode.GetStatus():
			c.JSON(http.StatusConflict, SetResponse(false, ConflictErrorCode, NULL()))
			c.Abort()
		default:
			c.JSON(http.StatusInternalServerError, SetResponse(false, ServerErrorCode, NULL()))
			c.Abort()

		}
	}
}
