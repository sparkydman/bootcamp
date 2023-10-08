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

// user providers
var userRepo = wire.NewSet(repository.NewUserRepository, wire.Bind(new(repository.IUserRepository), new(*repository.UserRepository)))

var userSvc = wire.NewSet(service.NewUserService, wire.Bind(new(service.IUserService), new(*service.UserService)))

var userCtrl = wire.NewSet(controller.NewUserController, wire.Bind(new(controller.IUserController), new(*controller.UserController)))

// bootcamp providers
var bcpRepo = wire.NewSet(repository.NewBootcampRepository, wire.Bind(new(repository.IBootcampRepository), new(*repository.BootcampRepository)))

var bcpSvc = wire.NewSet(service.NewBootcampService, wire.Bind(new(service.IBootcampService), new(*service.Bootcamp)))

var bcpCtrl = wire.NewSet(controller.NewBootcampController, wire.Bind(new(controller.IBootcampController), new(*controller.BootcampController)))

func InitApp(uri string) (*Initialize, error) {
	wire.Build(NewInitialize, db, userRepo, userSvc, userCtrl, bcpRepo, bcpSvc, bcpCtrl)
	return &Initialize{}, nil
}
