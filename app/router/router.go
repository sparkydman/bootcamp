package router

import (
	"bootcamp-api/app/middleware"
	"bootcamp-api/dependencies"

	"github.com/gin-gonic/gin"
)

func Init(app *dependencies.Initialize) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		//user routes
		user := api.Group("/users")
		user.POST("/create", app.UserCtrl.CreateUser)
		user.POST("/login", app.UserCtrl.LoginUser)
		user.GET("/", app.UserCtrl.GetUsers)
		user.GET("/me", middleware.AuthenticateUser(), app.UserCtrl.GetLoggedInUser)
		user.GET("/token", app.UserCtrl.GetToken)
		user.GET("/logout", middleware.AuthenticateUser(), app.UserCtrl.Logout)
		user.GET("/:id", app.UserCtrl.GetUserById)
		user.PUT("/:userid", middleware.AuthenticateUser(), app.UserCtrl.UpdateUser)
		user.DELETE("/:userid", middleware.AuthenticateUser(), app.UserCtrl.DeleteUser)

		//bootcamp routes
		bootcamp := api.Group("/bootcamps")
		bootcamp.POST("/", middleware.AuthenticateUser(), middleware.AuthorizeUser("admin", "publisher"), app.BootcampCtrl.AddBootcamp)
		bootcamp.GET("/creator-bootcamps/:creator", app.BootcampCtrl.GetBootcampsByCreator)
	}

	return router
}
