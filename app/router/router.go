package router

import (
	"bootcamp-api/app/controller"
	"bootcamp-api/app/middleware"

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
		user.PUT("/:userid", middleware.AuthenticateUser(), ctrl.UpdateUser)
		user.DELETE("/:userid", middleware.AuthenticateUser(), ctrl.DeleteUser)
	}

	return router
}
