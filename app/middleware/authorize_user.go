package middleware

import (
	"bootcamp-api/app/model/dao"
	"bootcamp-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeUser(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("LoggedInUser").(dao.User)
		isAllowed := false
		for _, r := range roles {
			if user.Role == r {
				isAllowed = true
			}
		}
		if isAllowed {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, utils.SetResponse(false, utils.ForbiddenErrorCode, utils.NULL()))
			c.Abort()
		}
	}
}
