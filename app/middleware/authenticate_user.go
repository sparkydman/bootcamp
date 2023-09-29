package middleware

import (
	"bootcamp-api/app/model/dao"
	"bootcamp-api/config"
	"bootcamp-api/utils"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := c.Request.Header
		authHeders := headers["Authorization"]
		if len(authHeders) != 0 {
			token := authHeders[0]
			if token != "" || strings.HasPrefix(token, "Bearer ") {
				token, found := strings.CutPrefix(token, "Bearer ")
				if found {
					user, err := config.VerifyToken(token, []byte(os.Getenv("ACCESS_TOKEN_KEY")), dao.User{})
					if err == nil {
						c.Set("LoggedInUser", user)
						c.Next()
					} else {
						utils.UnAuthorizedResponse(c)
					}
				}
			} else {
				utils.UnAuthorizedResponse(c)
			}
		} else {
			utils.UnAuthorizedResponse(c)
		}
	}
}
