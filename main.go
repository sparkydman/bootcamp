package main

import (
	"bootcamp-api/app/controller"
	"bootcamp-api/app/repository"
	"bootcamp-api/app/router"
	"bootcamp-api/app/service"
	"bootcamp-api/config"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db := config.Db{
		Uri: os.Getenv("MONGO_URI"),
	}
	if err := db.Connect(); err != nil {
		panic(err)
	}
	if db.IsConnected() {
		fmt.Println("Connection to database established")
	}
	defer db.Disconnect()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	app := router.Init(userController)

	app.Run(":8800")
}
