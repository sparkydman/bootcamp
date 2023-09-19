package router

import (
	"bootcamp-api/app/controller"

	"github.com/gin-gonic/gin"
)

func Init(ctrl controller.IUserController) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		user := api.Group("/users")
		user.POST("/create", ctrl.CreateUser)
		user.POST("/login", ctrl.LoginUser)
		user.GET("/", ctrl.GetUsers)
		user.GET("/:id", ctrl.GetUserById)
		user.PUT("/:userid", ctrl.UpdateUser)
		user.DELETE("/:userid", ctrl.DeleteUser)
	}

	return router
}
