package main

import (
	"bootcamp-api/app/router"
	"bootcamp-api/config"
	"bootcamp-api/dependencies"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func init() {
	config.InitLog()
}

func main() {
	dep, err := dependencies.InitApp(os.Getenv("MONGO_URI"))
	if err != nil {
		panic(err)
	}

	if dep.Db.IsConnected() {
		log.Println("Connection to database established")
	}

	defer dep.Db.Disconnect()

	app := router.Init(dep)

	app.Run(":8800")
}
