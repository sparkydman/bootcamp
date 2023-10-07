// go:build wireinject
//go:build wireinject
// +build wireinject

package dependencies

import (
	"bootcamp-api/app/controller"
	"bootcamp-api/app/repository"
	"bootcamp-api/app/service"
	"bootcamp-api/config"

	"github.com/google/wire"
)

var db = wire.NewSet(config.NewConnection)

var userRepo = wire.NewSet(repository.NewUserRepository, wire.Bind(new(repository.IUserRepository), new(*repository.UserRepository)))

var userSvc = wire.NewSet(service.NewUserService, wire.Bind(new(service.IUserService), new(*service.UserService)))

var userCtrl = wire.NewSet(controller.NewUserController, wire.Bind(new(controller.IUserController), new(*controller.UserController)))

func InitApp(uri string) (*Initialize, error) {
	wire.Build(NewInitialize, db, userRepo, userSvc, userCtrl)
	return &Initialize{}, nil
}
