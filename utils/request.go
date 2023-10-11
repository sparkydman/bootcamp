package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindOptions(c *gin.Context) []*options.FindOptions {
	querySort := c.Query("sort")
	querySelect := c.Query("select")
	opts := []*options.FindOptions{}

	if querySort != "" {
		opts = append(opts, options.Find().SetSort(bson.D{{querySort, 1}}))
	}
	if querySelect != "" {
		fields := strings.Split(querySelect, ",")
		var s bson.D
		for _, field := range fields {
			s = append(s, bson.E{Key: field, Value: 1})
		}
		opts = append(opts, options.Find().SetProjection(s))
	}
	return opts
}
