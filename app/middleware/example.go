package middleware

import "github.com/gin-gonic/gin"

func ExampleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("RequestUrl", c.Request.URL.RequestURI())

		c.Next()
	}
}
